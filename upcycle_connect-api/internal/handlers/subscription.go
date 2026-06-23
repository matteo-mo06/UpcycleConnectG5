package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/stripe/stripe-go/v85"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
)

func GetSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	plans, err := db.GetSubscriptionPlans()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	_ = json.NewEncoder(w).Encode(plans)
}

func GetMySubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	sub, err := db.GetUserActiveSubscription(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if sub == nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"active": false})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"active": true, "subscription": sub})
}

func CreateSubscriptionCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var body struct {
		PlanID string `json:"plan_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.PlanID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "plan_id requis"})
		return
	}

	plan, err := db.GetSubscriptionPlanByID(body.PlanID)
	if err != nil || !plan.IsActive {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "plan introuvable ou inactif"})
		return
	}
	if plan.StripePriceID == nil || *plan.StripePriceID == "" {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "paiement non disponible pour ce plan"})
		return
	}

	existing, _ := db.GetUserActiveSubscription(userID)
	if existing != nil {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous avez déjà un abonnement actif"})
		return
	}

	user, err := db.GetUserById(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	client := stripe.NewClient(config.StripeSecretKey())

	sessionParams := &stripe.CheckoutSessionCreateParams{
		Mode: stripe.String("subscription"),
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				Price:    stripe.String(*plan.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(config.FrontendURL() + "/abonnement?success=1"),
		CancelURL:  stripe.String(config.FrontendURL() + "/abonnement?cancelled=1"),
		CustomerEmail: stripe.String(user.Email),
		Metadata: map[string]string{
			"user_id": userID,
			"plan_id": plan.IdPlan,
		},
	}

	existingCustomerID := db.GetUserStripeCustomerID(userID)
	if existingCustomerID != "" {
		sessionParams.Customer = stripe.String(existingCustomerID)
		sessionParams.CustomerEmail = nil
	}

	session, err := client.V1CheckoutSessions.Create(context.Background(), sessionParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création du paiement"})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"url": session.URL})
}

func GetSubscriptionPortal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	customerID := db.GetUserStripeCustomerID(userID)
	if customerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "aucun abonnement trouvé"})
		return
	}

	client := stripe.NewClient(config.StripeSecretKey())
	session, err := client.V1BillingPortalSessions.Create(context.Background(), &stripe.BillingPortalSessionCreateParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(config.FrontendURL() + "/abonnement"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de l'ouverture du portail"})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"url": session.URL})
}


func GetAdminSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	plans, err := db.GetSubscriptionPlans()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	_ = json.NewEncoder(w).Encode(plans)
}

func UpdateAdminSubscriptionPlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	var req models.UpdateSubscriptionPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "données invalides"})
		return
	}

	current, err := db.GetSubscriptionPlanByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "plan introuvable"})
		return
	}

	client := stripe.NewClient(config.StripeSecretKey())

	newPriceID := current.StripePriceID
	newProductID := current.StripeProductID

	if current.StripeProductID == nil {
		product, err := client.V1Products.Create(context.Background(), &stripe.ProductCreateParams{
			Name: stripe.String(req.Name),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur Stripe lors de la création du produit"})
			return
		}
		newProductID = stripe.String(product.ID)
		price, err := client.V1Prices.Create(context.Background(), &stripe.PriceCreateParams{
			Product:    stripe.String(product.ID),
			Currency:   stripe.String("eur"),
			UnitAmount: stripe.Int64(int64(req.PriceCents)),
			Recurring: &stripe.PriceCreateRecurringParams{
				Interval:      stripe.String(req.IntervalUnit),
				IntervalCount: stripe.Int64(int64(req.IntervalCount)),
			},
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur Stripe lors de la création du prix"})
			return
		}
		newPriceID = stripe.String(price.ID)
	} else {
		if req.Name != current.Name {
			_, err = client.V1Products.Update(context.Background(), *current.StripeProductID, &stripe.ProductUpdateParams{
				Name: stripe.String(req.Name),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur Stripe lors de la mise à jour du produit"})
				return
			}
		}

		priceChanged := req.PriceCents != current.PriceCents || req.IntervalUnit != current.IntervalUnit || req.IntervalCount != current.IntervalCount
		if priceChanged {
			if current.StripePriceID != nil {
				_, _ = client.V1Prices.Update(context.Background(), *current.StripePriceID, &stripe.PriceUpdateParams{
					Active: stripe.Bool(false),
				})
			}
			if req.IsActive {
				price, err := client.V1Prices.Create(context.Background(), &stripe.PriceCreateParams{
					Product:    stripe.String(*current.StripeProductID),
					Currency:   stripe.String("eur"),
					UnitAmount: stripe.Int64(int64(req.PriceCents)),
					Recurring: &stripe.PriceCreateRecurringParams{
						Interval:      stripe.String(req.IntervalUnit),
						IntervalCount: stripe.Int64(int64(req.IntervalCount)),
					},
				})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur Stripe lors de la création du nouveau prix"})
					return
				}
				newPriceID = stripe.String(price.ID)
			} else {
				newPriceID = nil
			}
		}

		if !priceChanged && !req.IsActive && current.IsActive && current.StripePriceID != nil {
			_, _ = client.V1Prices.Update(context.Background(), *current.StripePriceID, &stripe.PriceUpdateParams{
				Active: stripe.Bool(false),
			})
			newPriceID = nil
		}
	}

	plan := models.SubscriptionPlan{
		IdPlan:          id,
		Name:            req.Name,
		PriceCents:      req.PriceCents,
		IntervalUnit:    req.IntervalUnit,
		IntervalCount:   req.IntervalCount,
		IsActive:        req.IsActive,
		StripePriceID:   newPriceID,
		StripeProductID: newProductID,
	}

	if err := db.UpdateSubscriptionPlan(plan); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la mise à jour"})
		return
	}

	updated, _ := db.GetSubscriptionPlanByID(id)
	_ = json.NewEncoder(w).Encode(updated)
}


func GetAdminSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := db.GetAllSubscriptions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if list == nil {
		list = []models.UserSubscriptionSummary{}
	}
	_ = json.NewEncoder(w).Encode(list)
}

func GetMySubscriptionInvoicePDF(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	invoiceID := r.PathValue("id")
	if invoiceID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "id requis"})
		return
	}

	customerID := db.GetUserStripeCustomerID(userID)
	if customerID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "aucun compte de facturation"})
		return
	}

	user, err := db.GetUserById(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	client := stripe.NewClient(config.StripeSecretKey())
	inv, err := client.V1Invoices.Retrieve(context.Background(), invoiceID, nil)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "facture introuvable"})
		return
	}

	if inv.Customer == nil || inv.Customer.ID != customerID {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}

	pdfBytes, err := buildSubscriptionInvoicePDF(invoiceID, inv, user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur génération PDF"})
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="facture-%s.pdf"`, invoiceID))
	_, _ = w.Write(pdfBytes)
}

func buildSubscriptionInvoicePDF(invoiceID string, inv *stripe.Invoice, user models.User) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 15, 20)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	totalPrice := fmt.Sprintf("%.2f EUR", float64(inv.AmountPaid)/100)
	date := time.Unix(inv.Created, 0).Format("02/01/2006 15:04")

	numDisplay := inv.Number
	if numDisplay == "" {
		numDisplay = invoiceID
	}
	if len(numDisplay) > 26 {
		numDisplay = numDisplay[:26]
	}

	var planDesc, periodStr string
	if inv.Lines != nil && len(inv.Lines.Data) > 0 {
		line := inv.Lines.Data[0]
		planDesc = line.Description
		if line.Period != nil {
			start := time.Unix(line.Period.Start, 0).Format("02/01/2006")
			end := time.Unix(line.Period.End, 0).Format("02/01/2006")
			periodStr = start + " au " + end
		}
	}
	planDesc = strings.ReplaceAll(planDesc, "/ month", "/mois")
	planDesc = strings.ReplaceAll(planDesc, "/month", "/mois")
	planDesc = strings.ReplaceAll(planDesc, "/ year", "/an")
	planDesc = strings.ReplaceAll(planDesc, "/year", "/an")
	if planDesc == "" {
		planDesc = "Abonnement"
	}
	if len([]rune(planDesc)) > 42 {
		planDesc = string([]rune(planDesc)[:42]) + "..."
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
	pdf.Cell(70, 6, tr(user.FirstName+" "+user.LastName))
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetXY(115, 68)
	pdf.Cell(70, 5, tr("E-mail : "+user.Email))

	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, 92, 190, 92)

	pdf.SetFillColor(220, 220, 220)
	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetXY(20, 96)
	pdf.CellFormat(90, 7, tr("Désignation"), "1", 0, "L", true, 0, "")
	pdf.CellFormat(45, 7, tr("Période"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 7, "Total TTC", "1", 0, "R", true, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetXY(20, 103)
	pdf.CellFormat(90, 8, tr(planDesc), "1", 0, "L", true, 0, "")
	pdf.CellFormat(45, 8, tr(periodStr), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, totalPrice, "1", 0, "R", true, 0, "")

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
	pdf.CellFormat(40, 6, "Sous-total", "LTR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, totalPrice, "LTR", 0, "R", false, 0, "")

	pdf.SetXY(115, 175)
	pdf.CellFormat(40, 6, "TVA (0%)", "LR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "0.00 EUR", "LR", 0, "R", false, 0, "")

	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetXY(115, 181)
	pdf.CellFormat(40, 7, "Total TTC", "1", 0, "L", false, 0, "")
	pdf.CellFormat(35, 7, totalPrice, "1", 0, "R", false, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(40, 40, 40)
	pdf.SetXY(20, 196)
	pdf.CellFormat(170, 7, tr("Règlement   Payé par carte bancaire le ")+date+" - Total : "+totalPrice, "1", 0, "L", false, 0, "")

	pdf.SetTextColor(120, 120, 120)
	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetXY(20, 268)
	pdf.CellFormat(170, 4, tr("UpcycleConnect - contact@upcycleconnect.fr - upcycleconnect.fr"), "", 2, "C", false, 0, "")
	pdf.SetX(20)
	pdf.CellFormat(170, 4, tr("Ce document tient lieu de reçu de paiement."), "", 0, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetMyInvoices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	customerID := db.GetUserStripeCustomerID(userID)
	if customerID == "" {
		_ = json.NewEncoder(w).Encode(map[string]any{"invoices": []any{}, "total": 0, "pages": 0})
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 5

	client := stripe.NewClient(config.StripeSecretKey())
	list := client.V1Invoices.List(context.Background(), &stripe.InvoiceListParams{
		ListParams: stripe.ListParams{Limit: stripe.Int64(100)},
		Customer:   stripe.String(customerID),
	})

	type invoiceItem struct {
		ID        string `json:"id"`
		Date      string `json:"date"`
		Amount    int64  `json:"amount_cents"`
		Currency  string `json:"currency"`
		Status    string `json:"status"`
		PDFURL    string `json:"pdf_url"`
		HostedURL string `json:"hosted_url"`
	}

	var all []invoiceItem
	for inv, err := range list.All(context.Background()) {
		if err != nil {
			break
		}
		item := invoiceItem{
			ID:       inv.ID,
			Date:     time.Unix(inv.Created, 0).UTC().Format("2006-01-02"),
			Amount:   inv.AmountPaid,
			Currency: string(inv.Currency),
			Status:   string(inv.Status),
		}
		if inv.InvoicePDF != "" {
			item.PDFURL = inv.InvoicePDF
		}
		if inv.HostedInvoiceURL != "" {
			item.HostedURL = inv.HostedInvoiceURL
		}
		all = append(all, item)
	}

	total := len(all)
	pages := (total + perPage - 1) / perPage
	start := (page - 1) * perPage
	end := start + perPage
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"invoices": all[start:end],
		"total":    total,
		"pages":    pages,
		"page":     page,
	})
}
