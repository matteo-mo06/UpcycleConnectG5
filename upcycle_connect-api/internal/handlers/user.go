package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
	"upcycle_connect-api/internal/utils"
)

func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	if err := db.DeleteUser(userID); err != nil {
		fmt.Println("DeleteMyAccount error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "échec de la suppression du compte"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "compte supprimé"})
}

func UpdateMyAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var body struct {
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}

	body.AvatarURL = strings.TrimSpace(body.AvatarURL)
	if body.AvatarURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "avatar_url requis"})
		return
	}

	if err := db.UpdateUserAvatar(userID, body.AvatarURL); err != nil {
		fmt.Println("UpdateMyAvatar error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "échec de la mise à jour de l'avatar"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "avatar mis à jour", "avatar_url": body.AvatarURL})
}

func UpdateMyPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var body struct {
		Current string `json:"current"`
		New     string `json:"new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}

	if body.Current == "" || body.New == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "mot de passe actuel et nouveau requis"})
		return
	}

	if len(body.New) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le nouveau mot de passe doit contenir au moins 8 caractères"})
		return
	}

	user, err := db.GetUserById(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "utilisateur non trouvé"})
			return
		}
		fmt.Println("UpdateMyPassword error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "échec de la récupération de l'utilisateur"})
		return
	}

	if err := utils.CheckPassword(body.Current, user.PasswordUser); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "mot de passe actuel incorrect"})
		return
	}

	hashed, err := utils.HashPassword(body.New)
	if err != nil {
		fmt.Println("UpdateMyPassword hash error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "échec du hachage du mot de passe"})
		return
	}

	if err := db.UpdateUserPassword(userID, hashed); err != nil {
		fmt.Println("UpdateMyPassword error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "échec de la mise à jour du mot de passe"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "mot de passe mis à jour"})
}

func UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var body struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}

	body.FirstName = strings.TrimSpace(body.FirstName)
	body.LastName = strings.TrimSpace(body.LastName)
	body.Email = strings.TrimSpace(body.Email)

	if body.FirstName == "" || body.LastName == "" || body.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "prénom, nom et email sont requis"})
		return
	}

	current, err := db.GetUserById(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "utilisateur non trouvé"})
			return
		}
		fmt.Println("UpdateMyProfile error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "échec de la récupération de l'utilisateur"})
		return
	}

	if body.Email != current.Email {
		_, err := db.GetUserByEmail(body.Email)
		if err == nil {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "cet email est déjà utilisé"})
			return
		}
		if err != sql.ErrNoRows {
			fmt.Println("UpdateMyProfile email check error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "échec de la vérification de l'email"})
			return
		}
	}

	if err := db.UpdateUserProfile(userID, body.FirstName, body.LastName, body.Email); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "cet email est déjà utilisé"})
			return
		}
		fmt.Println("UpdateMyProfile error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "échec de la mise à jour du profil"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "profil mis à jour"})
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	user, err := db.GetUserById(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("GetMe error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get user"})
		return
	}

	roles, _ := db.GetUserRoleNames(userID)
	if roles == nil {
		roles = []string{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user.ToLoginResponse(roles, []string{}))
}

func MarkTutorialDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	if err := db.MarkTutorialDone(userID); err != nil {
		fmt.Println("MarkTutorialDone error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update tutorial status"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "tutorial marked as done"})
}

func ResetTutorial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	if err := db.ResetTutorial(userID); err != nil {
		fmt.Println("ResetTutorial error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to reset tutorial"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "tutorial reset"})
}

func GetMyLimits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)

	if !slices.Contains(roles, config.RoleArtisan) {
		_ = json.NewEncoder(w).Encode(map[string]any{"is_premium": true})
		return
	}

	sub, _ := db.GetUserActiveSubscription(userID)
	if sub == nil {
		annCount, _ := db.GetUserAnnouncementCount(userID)
		projectCount, _ := db.CountUserProjects(userID)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"is_premium":    false,
			"announcements": map[string]int{"used": annCount, "max": 3},
			"projects":      map[string]int{"used": projectCount, "max": 2},
		})
		return
	}

	resp := map[string]any{"is_premium": true, "plan": sub.Plan}

	if sub.Plan.MaxAnnouncementsPerMonth != nil {
		used, _ := db.GetUserAnnouncementCountThisMonth(userID)
		resp["announcements"] = map[string]int{"used": used, "max": *sub.Plan.MaxAnnouncementsPerMonth}
	}
	if sub.Plan.MaxProjectsPerMonth != nil {
		used, _ := db.CountUserProjectsThisMonth(userID)
		resp["projects"] = map[string]int{"used": used, "max": *sub.Plan.MaxProjectsPerMonth}
	}
	if sub.Plan.MaxFeaturesPerMonth != nil {
		used, _ := db.GetUserFeatureRequestCountThisMonth(userID)
		resp["features"] = map[string]int{"used": used, "max": *sub.Plan.MaxFeaturesPerMonth}
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func GetMyStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	stats, err := db.GetUserStats(userID)
	if err != nil {
		fmt.Println("GetMyStats error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get stats"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}

func UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	allowed := map[string]bool{"active": true, "suspended": true, "blacklisted": true}
	if !allowed[body.Status] {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "status invalide (active | suspended | blacklisted)"})
		return
	}

	if err := db.UpdateUserStatus(id, body.Status); err != nil {
		fmt.Println("UpdateUserStatus error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update status"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "statut mis à jour"})
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	user, err := db.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("GetUserById error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user.ToResponse())
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	role := r.URL.Query().Get("role")
	page, limit, offset := parsePage(r, 50)

	users, total, err := db.GetAllUsersWithRoles(search, status, role, limit, offset)
	if err != nil {
		fmt.Println("GetUsers error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch users"})
		return
	}

	if users == nil {
		users = []models.UserListItem{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(users, total, page, limit))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	_, err := db.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("DeleteUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find user"})
		return
	}

	err = db.DeleteUser(id)
	if err != nil {
		fmt.Println("DeleteUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "user deleted successfully"})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	_, err := db.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("UpdateUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find user"})
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	user.IdUser = id

	err = db.UpdateUser(user)
	if err != nil {
		fmt.Println("UpdateUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "user updated successfully", "user": user.ToResponse()})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	user.Email = strings.TrimSpace(user.Email)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	err = utils.ValidateUserCreation(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(user.PasswordUser)
	if err != nil {
		fmt.Println("CreateUser hash error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to hash password"})
		return
	}

	user.IdUser = uuid.New().String()
	user.PasswordUser = hashedPassword

	profileToRole := map[string]string{
		"particulier":   "user",
		"artisan":       "artisan",
		"professionnel": "artisan",
	}
	profile := strings.TrimSpace(user.Profile)
	roleName, ok := profileToRole[profile]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "profil invalide (particulier, artisan, professionnel)"})
		return
	}

	err = db.CreateUser(user)
	if err != nil {
		fmt.Println("CreateUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create user"})
		return
	}

	role, err := db.GetRoleByName(roleName)
	if err != nil {
		fmt.Println("CreateUser GetRoleByName error:", err)
	} else {
		if err = db.AddUserRole(user.IdUser, role.IdRole); err != nil {
			fmt.Println("CreateUser AddUserRole error:", err)
		}
	}

	if verifToken, err := utils.GenerateEmailVerificationToken(user.IdUser); err != nil {
		fmt.Println("CreateUser GenerateEmailVerificationToken error:", err)
	} else {
		link := config.PlatformURL() + "/api/auth/verify-email?token=" + verifToken
		if err := utils.SendVerificationEmail(user.Email, link); err != nil {
			fmt.Println("CreateUser SendVerificationEmail error:", err)
		}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "user created successfully, un email de vérification a été envoyé",
		"user":    user.ToResponse(),
	})
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(verifyEmailPage("Lien invalide", "Le lien de vérification est invalide ou incomplet.")))
		return
	}

	userID, err := utils.ParseEmailVerificationToken(token)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(verifyEmailPage("Lien invalide ou expiré", "Ce lien de vérification n'est plus valide. Demandez-en un nouveau depuis la page de connexion.")))
		return
	}

	if err := db.SetEmailVerified(userID); err != nil {
		fmt.Println("VerifyEmail error:", err)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(verifyEmailPage("Erreur", "Une erreur est survenue, réessayez plus tard.")))
		return
	}

	http.Redirect(w, r, config.FrontendURL()+"/login?email_verified=1", http.StatusFound)
}

func verifyEmailPage(title, message string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr"><head><meta charset="utf-8"><title>%s</title>
<style>body{font-family:sans-serif;max-width:480px;margin:80px auto;text-align:center;color:#222}h1{color:#2e7d32}</style>
</head><body><h1>%s</h1><p>%s</p></body></html>`, title, title, message)
}

func ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	email := strings.TrimSpace(payload.Email)
	if user, err := db.GetUserByEmail(email); err == nil && !user.EmailVerified {
		if verifToken, err := utils.GenerateEmailVerificationToken(user.IdUser); err == nil {
			link := config.PlatformURL() + "/api/auth/verify-email?token=" + verifToken
			if err := utils.SendVerificationEmail(user.Email, link); err != nil {
				fmt.Println("ResendVerificationEmail SendVerificationEmail error:", err)
			}
		}
	}

	// Reponse identique dans tous les cas (email inconnu ou deja verifie) pour eviter l'enumeration de comptes
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "si un compte existe avec cet email, un lien de vérification a été envoyé"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	payload.Email = strings.TrimSpace(payload.Email)

	if payload.Email == "" || payload.PasswordUser == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "email and password are required"})
		return
	}

	user, err := db.GetUserByEmail(payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
			return
		}
		fmt.Println("Login error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch user"})
		return
	}

	err = utils.CheckPassword(payload.PasswordUser, user.PasswordUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	if !user.EmailVerified {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "email non vérifié, merci de consulter votre boîte mail"})
		return
	}

	if user.Status == "blacklisted" {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "compte banni"})
		return
	}

	if user.Status == "suspended" {
		expired, err := db.IsSuspensionExpired(user.IdUser)
		if err != nil {
			fmt.Println("Login IsSuspensionExpired error:", err)
		}
		if expired {
			_ = db.UpdateUserStatus(user.IdUser, "active")
		} else {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "compte suspendu"})
			return
		}
	}

	roles, _ := db.GetUserRoleNames(user.IdUser)
	if roles == nil {
		roles = []string{}
	}

	permissions, _ := db.GetUserPermissions(user.IdUser)
	if permissions == nil {
		permissions = []string{}
	}

	token, err := utils.GenerateToken(user.IdUser, roles, permissions)
	if err != nil {
		fmt.Println("Login GenerateToken error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "login successful",
		"token":   token,
		"user":    user.ToLoginResponse(roles, permissions),
	})
}
