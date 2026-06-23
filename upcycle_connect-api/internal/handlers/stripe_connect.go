package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stripe/stripe-go/v85"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
)

func GetStripeConnectStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	accountID := db.GetStripeAccountID(userID)
	if accountID == "" {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"connected":       false,
			"charges_enabled": false,
		})
		return
	}

	client := stripe.NewClient(config.StripeSecretKey())
	account, err := client.V1Accounts.GetByID(context.Background(), accountID, nil)
	if err != nil {
		fmt.Println("GetStripeConnectStatus error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la récupération du compte"})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"connected":       true,
		"charges_enabled": account.ChargesEnabled,
		"account_id":      accountID,
	})
}

func CreateStripeConnectOnboarding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	client := stripe.NewClient(config.StripeSecretKey())

	accountID := db.GetStripeAccountID(userID)

	if accountID == "" {
		account, err := client.V1Accounts.Create(context.Background(), &stripe.AccountCreateParams{
			BusinessType: stripe.String("individual"),
			Country:      stripe.String("FR"),
			BusinessProfile: &stripe.AccountCreateBusinessProfileParams{
				URL: stripe.String(config.PlatformURL()),
			},
			Controller: &stripe.AccountCreateControllerParams{
				StripeDashboard: &stripe.AccountCreateControllerStripeDashboardParams{
					Type: stripe.String("express"),
				},
				Fees: &stripe.AccountCreateControllerFeesParams{
					Payer: stripe.String("application"),
				},
				Losses: &stripe.AccountCreateControllerLossesParams{
					Payments: stripe.String("application"),
				},
			},
		})
		if err != nil {
			fmt.Println("CreateStripeAccount error:", err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "stripe_connect_unavailable"})
			return
		}
		accountID = account.ID
		if err := db.SaveStripeAccountID(userID, accountID); err != nil {
			fmt.Println("SaveStripeAccountID error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la sauvegarde du compte"})
			return
		}
	}

	frontendURL := config.FrontendURL()
	link, err := client.V1AccountLinks.Create(context.Background(), &stripe.AccountLinkCreateParams{
		Account:    stripe.String(accountID),
		RefreshURL: stripe.String(frontendURL + "/stripe/refresh"),
		ReturnURL:  stripe.String(frontendURL + "/stripe/return"),
		Type:       stripe.String("account_onboarding"),
	})
	if err != nil {
		fmt.Println("CreateAccountLink error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création du lien"})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"url": link.URL})
}
