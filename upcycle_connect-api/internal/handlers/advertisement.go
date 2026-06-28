package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v85"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
	"upcycle_connect-api/internal/utils"
)

func GetActiveAdvertisements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := db.GetActiveAdvertisements()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if list == nil {
		list = []models.Advertisement{}
	}
	_ = json.NewEncoder(w).Encode(list)
}

func CreateAdvertisement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	if !db.IsUserPremium(userID) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "abonnement premium requis pour publier une publicité"})
		return
	}

	var body struct {
		Title    string  `json:"title"`
		ImageURL string  `json:"image_url"`
		LinkURL  *string `json:"link_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" || body.ImageURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "titre et image requis"})
		return
	}

	priceCents, _ := db.GetAdvertisementPriceCents()

	ad := models.Advertisement{
		IdAdvertisement: uuid.New().String(),
		IdUser:          userID,
		Title:           body.Title,
		ImageURL:        body.ImageURL,
		LinkURL:         body.LinkURL,
		PriceCents:      priceCents,
	}

	if err := db.CreateAdvertisement(ad); err != nil {
		fmt.Println("CreateAdvertisement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ad)
}

func GetMyAdvertisements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	list, err := db.GetUserAdvertisements(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if list == nil {
		list = []models.Advertisement{}
	}
	_ = json.NewEncoder(w).Encode(list)
}

func CreateAdvertisementCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	ad, err := db.GetAdvertisementByID(id)
	if err != nil || ad == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "publicité introuvable"})
		return
	}
	if ad.IdUser != userID {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}
	if ad.State != "approved" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "la publicité doit être approuvée avant le paiement"})
		return
	}

	client := stripe.NewClient(config.StripeSecretKey())
	session, err := client.V1CheckoutSessions.Create(context.Background(), &stripe.CheckoutSessionCreateParams{
		Mode: stripe.String("payment"),
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
					Currency:   stripe.String("eur"),
					UnitAmount: stripe.Int64(int64(ad.PriceCents)),
					ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
						Name: stripe.String("Publicité UpcycleConnect - " + ad.Title),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(config.FrontendURL() + "/artisan/publicites?success=1"),
		CancelURL:  stripe.String(config.FrontendURL() + "/artisan/publicites"),
		Metadata: map[string]string{
			"advertisement_id": id,
			"user_id":          userID,
			"type":             "advertisement",
		},
	})
	if err != nil {
		fmt.Println("CreateAdvertisementCheckout Stripe error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création du paiement"})
		return
	}

	_ = db.SaveAdvertisementCheckoutSession(id, session.ID)

	_ = json.NewEncoder(w).Encode(map[string]string{"url": session.URL})
}

func GetAdminAdvertisementStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats, err := db.GetAdvertisementStats()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	_ = json.NewEncoder(w).Encode(stats)
}

func GetAdminAdvertisements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	state := r.URL.Query().Get("state")
	search := r.URL.Query().Get("search")

	list, total, err := db.GetAllAdvertisementsPaginated(page, limit, state, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if list == nil {
		list = []models.Advertisement{}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"advertisements": list,
		"total":          total,
		"page":           page,
		"limit":          limit,
	})
}

func ApproveAdvertisement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	ad, err := db.GetAdvertisementByID(id)
	if err != nil || ad == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "publicité introuvable"})
		return
	}
	if err := db.ApproveAdvertisement(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	go utils.SendPushNotification(
		db.GetOnesignalPlayerID(ad.IdUser),
		"Publicité approuvée",
		"Votre publicité \""+ad.Title+"\" a été approuvée. Vous pouvez maintenant procéder au paiement.",
	)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "approuvée"})
}

func RejectAdvertisement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	var body struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	ad, err := db.GetAdvertisementByID(id)
	if err != nil || ad == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "publicité introuvable"})
		return
	}
	if err := db.RejectAdvertisement(id, body.Reason); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	go utils.SendPushNotification(
		db.GetOnesignalPlayerID(ad.IdUser),
		"Publicité refusée",
		"Votre publicité \""+ad.Title+"\" n'a pas été approuvée.",
	)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "refusée"})
}

func DeactivateAdvertisement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	if err := db.DeactivateAdvertisement(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "désactivée"})
}

func ReactivateAdvertisement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	if err := db.ReactivateAdvertisement(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "réactivée"})
}

func GetAdvertisementPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	price, err := db.GetAdvertisementPriceCents()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]int{"price_cents": price})
}

func UpdateAdvertisementPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body struct {
		PriceCents int `json:"price_cents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.PriceCents < 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "prix invalide"})
		return
	}
	if err := db.SetAdvertisementPriceCents(body.PriceCents); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]int{"price_cents": body.PriceCents})
}
