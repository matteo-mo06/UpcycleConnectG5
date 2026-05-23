package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
)

func GetForumTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, limit, offset := parsePage(r, 20)
	topics, total, err := db.GetTopics(r.URL.Query().Get("search"), limit, offset)
	if err != nil {
		fmt.Println("GetForumTopics error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch topics"})
		return
	}

	if topics == nil {
		topics = []models.Topic{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(topics, total, page, limit))
}

func GetForumTopic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing topic id"})
		return
	}

	topic, err := db.GetTopicByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "topic not found"})
			return
		}
		fmt.Println("GetForumTopic error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get topic"})
		return
	}

	if topic.Posts == nil {
		topic.Posts = []models.Post{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(topic)
}

func CreateForumTopic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var req models.CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if req.Title == "" || req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le titre et le message sont requis"})
		return
	}

	topicID, err := db.CreateTopic(userID, req.Title, req.Body)
	if err != nil {
		fmt.Println("CreateForumTopic error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create topic"})
		return
	}

	topic, err := db.GetTopicByID(topicID)
	if err != nil {
		fmt.Println("CreateForumTopic fetch error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "topic created but failed to fetch"})
		return
	}

	if topic.Posts == nil {
		topic.Posts = []models.Post{}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(topic)
}

func AddForumPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	topicID := r.PathValue("id")

	if topicID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing topic id"})
		return
	}

	_, err := db.GetTopicByID(topicID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "topic not found"})
			return
		}
		fmt.Println("AddForumPost error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find topic"})
		return
	}

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le message est requis"})
		return
	}

	if err := db.CreatePost(topicID, userID, req.Body, req.IdParentPost); err != nil {
		fmt.Println("AddForumPost error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create post"})
		return
	}

	topic, err := db.GetTopicByID(topicID)
	if err != nil {
		fmt.Println("AddForumPost fetch error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "post created but failed to fetch topic"})
		return
	}

	if topic.Posts == nil {
		topic.Posts = []models.Post{}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(topic)
}

func UpdateForumTopic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	id := r.PathValue("id")

	authorID, err := db.GetTopicAuthor(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "topic introuvable"})
			return
		}
		fmt.Println("UpdateForumTopic GetTopicAuthor error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	if userID != authorID && !slices.Contains(roles, config.RoleAdmin) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}

	var req models.UpdateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}
	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le titre est requis"})
		return
	}

	if err := db.UpdateTopicTitle(id, req.Title); err != nil {
		fmt.Println("UpdateForumTopic error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la mise à jour"})
		return
	}

	topic, err := db.GetTopicByID(id)
	if err != nil {
		fmt.Println("UpdateForumTopic fetch error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "mise à jour effectuée mais erreur lors de la récupération"})
		return
	}
	if topic.Posts == nil {
		topic.Posts = []models.Post{}
	}

	_ = json.NewEncoder(w).Encode(topic)
}

func DeleteForumTopic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	perms, _ := r.Context().Value(middleware.ContextPermissions).([]string)
	id := r.PathValue("id")

	authorID, err := db.GetTopicAuthor(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "topic introuvable"})
			return
		}
		fmt.Println("DeleteForumTopic GetTopicAuthor error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	isAdmin := slices.Contains(roles, config.RoleAdmin)
	canModerate := slices.Contains(perms, "moderate_forum")
	if userID != authorID && !isAdmin && !canModerate {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}

	if err := db.SoftDeleteContent("topic", id); err != nil {
		fmt.Println("DeleteForumTopic error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la suppression"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateForumPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	topicID := r.PathValue("id")
	postID := r.PathValue("post_id")

	authorID, err := db.GetPostAuthor(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "post introuvable"})
			return
		}
		fmt.Println("UpdateForumPost GetPostAuthor error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	if userID != authorID && !slices.Contains(roles, config.RoleAdmin) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}

	var req models.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}
	if req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le message est requis"})
		return
	}

	if err := db.UpdatePostBody(postID, req.Body); err != nil {
		fmt.Println("UpdateForumPost error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la mise à jour"})
		return
	}

	topic, err := db.GetTopicByID(topicID)
	if err != nil {
		fmt.Println("UpdateForumPost fetch error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "mise à jour effectuée mais erreur lors de la récupération"})
		return
	}
	if topic.Posts == nil {
		topic.Posts = []models.Post{}
	}

	_ = json.NewEncoder(w).Encode(topic)
}

func DeleteForumPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	perms, _ := r.Context().Value(middleware.ContextPermissions).([]string)
	topicID := r.PathValue("id")
	postID := r.PathValue("post_id")

	authorID, err := db.GetPostAuthor(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "post introuvable"})
			return
		}
		fmt.Println("DeleteForumPost GetPostAuthor error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	isAdmin := slices.Contains(roles, config.RoleAdmin)
	canModerate := slices.Contains(perms, "moderate_forum")
	if userID != authorID && !isAdmin && !canModerate {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}

	if err := db.SoftDeleteContent("post", postID); err != nil {
		fmt.Println("DeleteForumPost error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la suppression"})
		return
	}

	topic, err := db.GetTopicByID(topicID)
	if err != nil {
		fmt.Println("DeleteForumPost fetch error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "suppression effectuée mais erreur lors de la récupération"})
		return
	}
	if topic.Posts == nil {
		topic.Posts = []models.Post{}
	}

	_ = json.NewEncoder(w).Encode(topic)
}
