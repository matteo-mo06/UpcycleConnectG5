package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v85"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
	"upcycle_connect-api/internal/utils"
)

func CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	announcementID := r.PathValue("id")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	if slices.Contains(roles, config.RoleArtisan) && !db.IsUserPremium(userID) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "abonnement premium requis pour acheter des annonces"})
		return
	}

	ann, err := db.GetAnnouncementById(announcementID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "annonce introuvable"})
		return
	}

	if ann.StateAnnouncement != "Active" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette annonce n'est plus disponible"})
		return
	}

	if ann.Price <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette annonce est gratuite, pas de paiement requis"})
		return
	}

	sellerID, err := db.GetAnnouncementOwnerID(announcementID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	sellerStripeAccountID := db.GetStripeAccountID(sellerID)
	if sellerStripeAccountID == "" {
		w.WriteHeader(http.StatusPaymentRequired)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le vendeur n'a pas encore configuré son compte de paiement"})
		return
	}

	commissionRate, _ := db.GetCommissionRate()
	priceCents := int64(ann.Price * 100)
	commissionCents := int64(float64(priceCents) * commissionRate / 100.0)
	totalCents := priceCents + commissionCents

	client := stripe.NewClient(config.StripeSecretKey())

	params := &stripe.PaymentIntentCreateParams{
		Amount:                  stripe.Int64(totalCents),
		Currency:                stripe.String("eur"),
		ApplicationFeeAmount:    stripe.Int64(commissionCents),
		TransferData: &stripe.PaymentIntentCreateTransferDataParams{
			Destination: stripe.String(sellerStripeAccountID),
		},
		Metadata: map[string]string{
			"announcement_id":  announcementID,
			"buyer_id":         userID,
			"price_cents":      strconv.FormatInt(priceCents, 10),
			"commission_cents": strconv.FormatInt(commissionCents, 10),
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentCreateAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := client.V1PaymentIntents.Create(context.Background(), params)
	if err != nil {
		fmt.Println("Stripe CreatePaymentIntent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création du paiement"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"client_secret":    pi.ClientSecret,
		"announcement_id":  announcementID,
		"price_cents":      priceCents,
		"commission_cents": commissionCents,
		"total_cents":      totalCents,
		"commission_rate":  commissionRate,
	})
}

func CreateFormationPaymentIntent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	formationID := r.PathValue("id")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	formation, err := db.GetFormationById(formationID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation introuvable"})
		return
	}

	if formation.Status != "approved" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette formation n'est pas disponible"})
		return
	}

	if formation.Price == nil || *formation.Price <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette formation est gratuite"})
		return
	}

	priceCents := int64(*formation.Price * 100)

	client := stripe.NewClient(config.StripeSecretKey())

	params := &stripe.PaymentIntentCreateParams{
		Amount:   stripe.Int64(priceCents),
		Currency: stripe.String("eur"),
		Metadata: map[string]string{
			"formation_id": formationID,
			"buyer_id":     userID,
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentCreateAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := client.V1PaymentIntents.Create(context.Background(), params)
	if err != nil {
		fmt.Println("Stripe CreateFormationPaymentIntent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création du paiement"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"client_secret": pi.ClientSecret,
		"price_cents":   priceCents,
	})
}

func StripeWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sig := r.Header.Get("Stripe-Signature")
	webhookSecret := config.StripeWebhookSecret()

	var event stripe.Event

	if webhookSecret != "" {
		event, err = stripe.ConstructEvent(payload, sig, webhookSecret)
		if err != nil {
			fmt.Println("Stripe webhook signature error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		if err := json.Unmarshal(payload, &event); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if session.Mode == "subscription" {
			userID := session.Metadata["user_id"]
			planID := session.Metadata["plan_id"]
			stripeSubID := ""
			if session.Subscription != nil {
				stripeSubID = session.Subscription.ID
			}
			customerID := ""
			if session.Customer != nil {
				customerID = session.Customer.ID
			}
			if userID != "" && planID != "" && stripeSubID != "" && !db.SubscriptionExistsByStripeID(stripeSubID) {
				subID := uuid.New().String()
				if err := db.CreateSubscriptionRecord(subID, stripeSubID, customerID, planID, userID); err != nil {
					fmt.Println("CreateSubscriptionRecord error:", err)
				}
			}
		}
		w.WriteHeader(http.StatusOK)
		return

	case "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var endStr *string
		if sub.CancelAt > 0 {
			t := time.Unix(sub.CancelAt, 0).UTC().Format("2006-01-02 15:04:05")
			endStr = &t
		}
		cancelled := sub.Status == "canceled" || sub.CancelAtPeriodEnd
		_ = db.UpdateSubscriptionByStripeID(sub.ID, endStr, cancelled)
		w.WriteHeader(http.StatusOK)
		return

	case "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var endStr *string
		if sub.EndedAt > 0 {
			t := time.Unix(sub.EndedAt, 0).UTC().Format("2006-01-02 15:04:05")
			endStr = &t
		}
		_ = db.UpdateSubscriptionByStripeID(sub.ID, endStr, true)
		w.WriteHeader(http.StatusOK)
		return
	}

	if event.Type == stripe.EventTypePaymentIntentSucceeded {
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			fmt.Println("Stripe webhook unmarshal error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		announcementID := pi.Metadata["announcement_id"]
		formationID := pi.Metadata["formation_id"]
		buyerID := pi.Metadata["buyer_id"]

		if formationID != "" {
			if err := db.RegisterUserForFormation(buyerID, formationID); err != nil {
				fmt.Println("RegisterUserForFormation (webhook) error:", err)
			}
		}

		if announcementID != "" {
			if err := db.MarkAnnouncementSold(announcementID, buyerID); err != nil {
				fmt.Println("MarkAnnouncementSold error:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			commissionCents, _ := strconv.Atoi(pi.Metadata["commission_cents"])
			if err := db.StorePayment(pi.ID, announcementID, buyerID, string(pi.Currency), pi.ID, int(pi.Amount), commissionCents); err != nil {
				fmt.Println("StorePayment error:", err)
			}

			go generateAndStoreInvoice(pi.ID, announcementID, buyerID, pi.Amount, int64(commissionCents), pi.Created)

			if sellerID, err := db.GetAnnouncementOwnerID(announcementID); err == nil {
				go utils.SendPushNotification(db.GetOnesignalPlayerID(sellerID), "Annonce vendue", "Votre annonce a été achetée !")
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func generateAndStoreInvoice(paymentID, announcementID, buyerID string, amountCents, commissionCents, createdAt int64) {
	ann, err := db.GetAnnouncementById(announcementID)
	if err != nil {
		fmt.Println("generateInvoice: GetAnnouncementById error:", err)
		return
	}
	buyer, err := db.GetUserById(buyerID)
	if err != nil {
		fmt.Println("generateInvoice: GetUserById error:", err)
		return
	}

	dir := config.InvoicesDir()
	filename := paymentID + ".pdf"
	fullPath := filepath.Join(dir, filename)

	if err := buildInvoicePDF(fullPath, paymentID, ann, buyer, amountCents, commissionCents, createdAt); err != nil {
		fmt.Println("generateInvoice: buildInvoicePDF error:", err)
		return
	}

	if err := db.StoreInvoicePath(announcementID, filename); err != nil {
		fmt.Println("generateInvoice: StoreInvoicePath error:", err)
	}
}

func buildInvoicePDF(destPath, paymentID string, ann models.Announcement, buyer models.User, amountCents, commissionCents, createdAt int64) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 15, 20)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	priceCents := amountCents - commissionCents
	articlePrice := fmt.Sprintf("%.2f EUR", float64(priceCents)/100)
	commissionPrice := fmt.Sprintf("%.2f EUR", float64(commissionCents)/100)
	totalPrice := fmt.Sprintf("%.2f EUR", float64(amountCents)/100)
	date := time.Unix(createdAt, 0).Format("02/01/2006 15:04")

	numDisplay := paymentID
	if len(numDisplay) > 26 {
		numDisplay = numDisplay[:26]
	}

	pdf.SetTextColor(45, 106, 79)
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetXY(20, 15)
	pdf.Cell(80, 10, "UpcycleConnect")

	pdf.SetDrawColor(150, 150, 150)
	pdf.Rect(122, 10, 68, 38, "D")
	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(122, 12)
	pdf.CellFormat(68, 7, "Facture", "", 0, "C", false, 0, "")
	pdf.Line(122, 20, 190, 20)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetXY(124, 22)
	pdf.CellFormat(64, 5, tr("Numéro : ")+numDisplay, "", 2, "L", false, 0, "")
	pdf.SetX(124)
	pdf.CellFormat(64, 5, "Date : "+date, "", 2, "L", false, 0, "")
	pdf.SetX(124)
	pdf.CellFormat(64, 5, tr("Page n°1 sur 1"), "", 0, "L", false, 0, "")

	pdf.SetTextColor(60, 60, 60)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetXY(20, 56)
	pdf.Cell(85, 6, "UpcycleConnect")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetXY(20, 62)
	pdf.Cell(85, 5, "contact@upcycleconnect.fr")
	pdf.SetXY(20, 67)
	pdf.Cell(85, 5, "upcycleconnect.fr")

	pdf.SetDrawColor(150, 150, 150)
	pdf.Rect(110, 55, 80, 34, "D")
	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(118, 52, 16, 6, "F")
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.SetXY(121, 52)
	pdf.Cell(16, 6, "Client")
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(40, 40, 40)
	pdf.SetXY(115, 61)
	pdf.Cell(70, 6, tr(buyer.FirstName+" "+buyer.LastName))
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetXY(115, 68)
	pdf.Cell(70, 5, tr("E-mail : "+buyer.Email))

	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, 92, 190, 92)

	pdf.SetFillColor(220, 220, 220)
	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetXY(20, 96)
	pdf.CellFormat(110, 7, tr("Désignation"), "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 7, tr("Quantité"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 7, "Total TTC", "1", 0, "R", true, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetFillColor(255, 255, 255)
	title := ann.TitleAnnouncement
	if len([]rune(title)) > 58 {
		title = string([]rune(title)[:58]) + "..."
	}
	pdf.SetXY(20, 103)
	pdf.CellFormat(110, 8, tr(title), "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 8, "1", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, articlePrice, "1", 0, "R", true, 0, "")

	pdf.SetFont("Helvetica", "I", 7.5)
	pdf.SetTextColor(120, 120, 120)
	pdf.SetXY(22, 112)
	pdf.Cell(108, 4, tr("Vendeur : "+ann.AuthorName))

	pdf.SetTextColor(220, 90, 90)
	pdf.SetFont("Helvetica", "B", 72)
	pdf.TransformBegin()
	pdf.TransformRotate(38, 95, 145)
	pdf.SetXY(28, 133)
	pdf.Cell(135, 25, tr("Payé"))
	pdf.TransformEnd()

	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetXY(115, 169)
	pdf.CellFormat(40, 6, "Prix article", "LTR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, articlePrice, "LTR", 0, "R", false, 0, "")

	pdf.SetXY(115, 175)
	pdf.CellFormat(40, 6, tr("Commission plateforme"), "LR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, commissionPrice, "LR", 0, "R", false, 0, "")

	pdf.SetXY(115, 181)
	pdf.CellFormat(40, 6, "TVA (0%)", "LR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "0.00 EUR", "LR", 0, "R", false, 0, "")

	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetXY(115, 187)
	pdf.CellFormat(40, 7, "Total TTC", "1", 0, "L", false, 0, "")
	pdf.CellFormat(35, 7, totalPrice, "1", 0, "R", false, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(40, 40, 40)
	pdf.SetXY(20, 202)
	pdf.CellFormat(170, 7, tr("Règlement   Payé par carte bancaire le ")+date+" - Total : "+totalPrice, "1", 0, "L", false, 0, "")

	pdf.SetTextColor(120, 120, 120)
	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetXY(20, 268)
	pdf.CellFormat(170, 4, tr("UpcycleConnect - contact@upcycleconnect.fr - upcycleconnect.fr"), "", 2, "C", false, 0, "")
	pdf.SetX(20)
	pdf.CellFormat(170, 4, tr("Ce document tient lieu de reçu de paiement."), "", 0, "C", false, 0, "")

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return pdf.Output(f)
}
