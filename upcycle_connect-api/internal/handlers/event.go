package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, limit, offset := parsePage(r, 20)
	events, total, err := db.GetAllEvents(r.URL.Query().Get("search"), r.URL.Query().Get("status"), limit, offset)
	if err != nil {
		fmt.Println("GetEvents error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch events"})
		return
	}

	if events == nil {
		events = []models.Event{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(events, total, page, limit))
}

func GetEventById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing event id"})
		return
	}

	event, err := db.GetEventById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "event not found"})
			return
		}
		fmt.Println("GetEventById error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get event"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(event)
}

func createEvent(w http.ResponseWriter, r *http.Request, defaultStatus string) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var e models.Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if e.TitleEvent == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "title is required"})
		return
	}

	if e.DateEvent == nil || *e.DateEvent == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "date is required"})
		return
	}

	parsed, err := time.Parse("2006-01-02T15:04:05", *e.DateEvent)
	if err != nil {
		parsed, err = time.Parse("2006-01-02T15:04", *e.DateEvent)
	}
	if err != nil || parsed.Before(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "la date doit être dans le futur"})
		return
	}

	e.IdEvent = uuid.New().String()
	if e.IdCreator == nil && userID != "" {
		e.IdCreator = &userID
	}
	e.Status = defaultStatus

	err = db.CreateEvent(e)
	if err != nil {
		fmt.Println("CreateEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create event"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "event created successfully", "event": e})
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	createEvent(w, r, "pending")
}

func CreateEventAdmin(w http.ResponseWriter, r *http.Request) {
	createEvent(w, r, "approved")
}

func ApproveEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	if err := db.ApproveEvent(id); err != nil {
		fmt.Println("ApproveEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to approve event"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "événement approuvé"})
}

func RejectEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	var body struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if err := db.RejectEvent(id, body.Reason); err != nil {
		fmt.Println("RejectEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to reject event"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "événement rejeté"})
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing event id"})
		return
	}

	_, err := db.GetEventById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "event not found"})
			return
		}
		fmt.Println("UpdateEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find event"})
		return
	}

	var e models.Event
	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	e.IdEvent = id

	err = db.UpdateEvent(e)
	if err != nil {
		fmt.Println("UpdateEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update event"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "event updated successfully", "event": e})
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing event id"})
		return
	}

	_, err := db.GetEventById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "event not found"})
			return
		}
		fmt.Println("DeleteEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find event"})
		return
	}

	err = db.DeleteEvent(id)
	if err != nil {
		fmt.Println("DeleteEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete event"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "event deleted successfully"})
}

func DeleteMyEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	event, err := db.GetEventById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "event not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find event"})
		return
	}

	if event.IdCreator == nil || *event.IdCreator != userID {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous n'êtes pas créateur de cet événement"})
		return
	}

	if err := db.DeleteEvent(id); err != nil {
		fmt.Println("DeleteMyEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete event"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "événement supprimé"})
}

func UpdateMyEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	event, err := db.GetEventById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "événement introuvable"})
			return
		}
		fmt.Println("UpdateMyEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find event"})
		return
	}

	if event.IdCreator == nil || *event.IdCreator != userID {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous n'êtes pas créateur de cet événement"})
		return
	}

	if event.Status == "approved" {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "un événement approuvé ne peut pas être modifié directement"})
		return
	}

	var body models.Event
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Champs de contenu uniquement — Id_creator et Status préservés
	event.TitleEvent = body.TitleEvent
	event.DescriptionEvent = body.DescriptionEvent
	event.DateEvent = body.DateEvent
	event.LocationEvent = body.LocationEvent
	event.Capacity = body.Capacity
	event.PriceCents = body.PriceCents

	if err := db.UpdateEvent(event); err != nil {
		fmt.Println("UpdateMyEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update event"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "événement mis à jour"})
}

func GetPublicEventsForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	page, limit, offset := parsePage(r, 15)

	events, total, err := db.GetPublicEventsForUser(userID, r.URL.Query().Get("search"), r.URL.Query().Get("status"), limit, offset)
	if err != nil {
		fmt.Println("GetPublicEventsForUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch events"})
		return
	}

	if events == nil {
		events = []models.Event{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(events, total, page, limit))
}

func GetUserEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	events, err := db.GetUserRegisteredEvents(userID)
	if err != nil {
		fmt.Println("GetUserEvents error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch user events"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func GetMyCreatedEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	events, err := db.GetMyCreatedEvents(userID)
	if err != nil {
		fmt.Println("GetMyCreatedEvents error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch created events"})
		return
	}
	if events == nil {
		events = []models.Event{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func RegisterForEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing event id"})
		return
	}

	event, err := db.GetEventById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "event not found"})
			return
		}
		fmt.Println("RegisterForEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find event"})
		return
	}

	if event.Status != "approved" {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cet événement n'est pas disponible"})
		return
	}

	err = db.RegisterUserForEvent(userID, id)
	if err != nil {
		fmt.Println("RegisterForEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to register for event"})
		return
	}

	if err := db.AwardScore(userID, "event_registration", id); err != nil {
		fmt.Println("AwardScore event_registration error:", err)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "registered successfully"})
}

func UnregisterFromEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing event id"})
		return
	}

	err := db.UnregisterUserFromEvent(userID, id)
	if err != nil {
		fmt.Println("UnregisterFromEvent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to unregister from event"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "unregistered successfully"})
}
