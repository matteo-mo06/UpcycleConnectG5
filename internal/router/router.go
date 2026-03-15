package router

import (
	"net/http"

	"upcycle_connect-api/internal/handlers"
	"upcycle_connect-api/internal/middleware"
)

func admin(h http.HandlerFunc) http.Handler {
	return middleware.RequireAdmin(h)
}

func auth(h http.HandlerFunc) http.Handler {
	return middleware.RequireAuth(h)
}

func InitRoutes() {
	http.HandleFunc("GET /", handlers.HealthCheck)

	http.HandleFunc("POST /auth/login", handlers.Login)
	http.HandleFunc("POST /auth/register", handlers.CreateUser)

	http.Handle("POST /professional-request", auth(handlers.SubmitProfessionalRequest))

	http.Handle("GET /admin/users", admin(handlers.GetUsers))
	http.Handle("GET /admin/user/{id}", admin(handlers.GetUserById))
	http.Handle("PUT /admin/user/{id}", admin(handlers.UpdateUser))
	http.Handle("DELETE /admin/user/{id}", admin(handlers.DeleteUser))

	http.Handle("GET /admin/roles", admin(handlers.GetRoles))
	http.Handle("GET /admin/user/{id}/roles", admin(handlers.GetUserRoles))
	http.Handle("POST /admin/user/{id}/roles", admin(handlers.AddRoleToUser))
	http.Handle("DELETE /admin/user/{id}/roles/{role_id}", admin(handlers.RemoveRoleFromUser))

	http.Handle("GET /admin/professional-requests", admin(handlers.GetProfessionalRequests))
	http.Handle("PUT /admin/professional-requests/{id}/validate", admin(handlers.ValidateProfessionalRequest))
	http.Handle("PUT /admin/professional-requests/{id}/reject", admin(handlers.RejectProfessionalRequest))

	http.Handle("GET /admin/categories", admin(handlers.GetCategories))
	http.Handle("GET /admin/category/{id}", admin(handlers.GetCategoryById))
	http.Handle("POST /admin/categories", admin(handlers.CreateCategory))
	http.Handle("PUT /admin/category/{id}", admin(handlers.UpdateCategory))
	http.Handle("DELETE /admin/category/{id}", admin(handlers.DeleteCategory))

	http.Handle("GET /admin/announcements", admin(handlers.GetAnnouncements))
	http.Handle("GET /admin/announcement/{id}", admin(handlers.GetAnnouncementById))
	http.Handle("POST /admin/announcements", admin(handlers.CreateAnnouncement))
	http.Handle("PUT /admin/announcement/{id}", admin(handlers.UpdateAnnouncement))
	http.Handle("DELETE /admin/announcement/{id}", admin(handlers.DeleteAnnouncement))

	http.Handle("GET /admin/events", admin(handlers.GetEvents))
	http.Handle("GET /admin/event/{id}", admin(handlers.GetEventById))
	http.Handle("POST /admin/events", admin(handlers.CreateEvent))
	http.Handle("PUT /admin/event/{id}", admin(handlers.UpdateEvent))
	http.Handle("DELETE /admin/event/{id}", admin(handlers.DeleteEvent))
}
