package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/stripe/stripe-go/v85"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
)

func CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	announcementID := r.PathValue("id")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

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

	client := stripe.NewClient(config.StripeSecretKey())

	params := &stripe.PaymentIntentCreateParams{
		Amount:   stripe.Int64(int64(ann.Price * 100)),
		Currency: stripe.String("eur"),
		Metadata: map[string]string{
			"announcement_id": announcementID,
			"buyer_id":        userID,
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
	_ = json.NewEncoder(w).Encode(map[string]string{
		"client_secret":   pi.ClientSecret,
		"announcement_id": announcementID,
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

	if event.Type == stripe.EventTypePaymentIntentSucceeded {
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			fmt.Println("Stripe webhook unmarshal error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		announcementID := pi.Metadata["announcement_id"]
		buyerID := pi.Metadata["buyer_id"]

		if announcementID != "" {
			if err := db.MarkAnnouncementSold(announcementID, buyerID); err != nil {
				fmt.Println("MarkAnnouncementSold error:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
