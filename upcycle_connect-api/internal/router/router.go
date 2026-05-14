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

func perm(permission string, h http.HandlerFunc) http.Handler {
	return middleware.RequirePermission(permission)(h)
}

func InitRoutes() {
	http.HandleFunc("GET /{$}", handlers.HealthCheck)
	http.Handle("GET /sse", auth(handlers.SSEHandler))

	http.HandleFunc("POST /auth/login", handlers.Login)
	http.HandleFunc("POST /auth/register", handlers.CreateUser)

	http.Handle("GET /user/me", auth(handlers.GetMe))
	http.Handle("GET /user/stats", auth(handlers.GetMyStats))
	http.Handle("GET /user/announcements", auth(handlers.GetMyAnnouncements))
	http.Handle("GET /user/acquisitions", auth(handlers.GetMyAcquisitions))
	http.Handle("DELETE /user/announcement/{id}", auth(handlers.DeleteMyAnnouncement))

	http.Handle("POST /upload", auth(handlers.UploadFile))

	http.HandleFunc("GET /announcements", handlers.GetPublicAnnouncements)
	http.HandleFunc("GET /announcements/{id}", handlers.GetPublicAnnouncementById)
	http.Handle("POST /announcements", perm("create_announcement", handlers.CreateUserAnnouncement))
	http.Handle("GET /announcements/{id}/documents", auth(handlers.GetAnnouncementDocuments))
	http.Handle("POST /announcements/{id}/claim", perm("buy_announcement", handlers.ClaimAnnouncement))
	http.Handle("DELETE /announcements/{id}", perm("manage_announcements", handlers.DeleteAnnouncementWithPermission))

	http.Handle("GET /user/deposit-requests", auth(handlers.GetMyDepositRequests))
	http.Handle("POST /announcements/{id}/deposit-request", perm("request_locker_deposit", handlers.CreateDepositRequest))
	http.Handle("DELETE /announcements/{id}/deposit-request", perm("request_locker_deposit", handlers.CancelDepositRequest))
	http.Handle("POST /deposit-requests/{id}/validate", perm("validate_deposit", handlers.ValidateDepositRequest))
	http.Handle("POST /deposit-requests/{id}/reject", perm("validate_deposit", handlers.RejectDepositRequest))
	http.Handle("GET /deposit-requests/pending", perm("validate_deposit", handlers.GetPendingDepositRequests))

	http.Handle("POST /professional-request", auth(handlers.SubmitProfessionalRequest))

	http.Handle("GET /admin/stats", admin(handlers.GetStats))

	http.Handle("GET /admin/users", admin(handlers.GetUsers))
	http.Handle("GET /admin/user/{id}", admin(handlers.GetUserById))
	http.Handle("PUT /admin/user/{id}", admin(handlers.UpdateUser))
	http.Handle("PATCH /admin/user/{id}/status", admin(handlers.UpdateUserStatus))
	http.Handle("DELETE /admin/user/{id}", admin(handlers.DeleteUser))

	http.Handle("GET /admin/roles", admin(handlers.GetRoles))
	http.Handle("GET /admin/user/{id}/roles", admin(handlers.GetUserRoles))
	http.Handle("POST /admin/user/{id}/roles", admin(handlers.AddRoleToUser))
	http.Handle("DELETE /admin/user/{id}/roles/{role_id}", admin(handlers.RemoveRoleFromUser))

	http.Handle("GET /admin/permissions", admin(handlers.GetPermissions))
	http.Handle("GET /admin/role/{id}/permissions", admin(handlers.GetRolePermissions))
	http.Handle("POST /admin/role/{id}/permissions", admin(handlers.AddRolePermission))
	http.Handle("DELETE /admin/role/{id}/permissions/{permission_id}", admin(handlers.RemoveRolePermission))

	http.Handle("GET /admin/professional-requests", admin(handlers.GetProfessionalRequests))
	http.Handle("PUT /admin/professional-requests/{id}/validate", admin(handlers.ValidateProfessionalRequest))
	http.Handle("PUT /admin/professional-requests/{id}/reject", admin(handlers.RejectProfessionalRequest))

	http.Handle("GET /admin/categories", admin(handlers.GetCategories))
	http.Handle("GET /admin/category/{id}", admin(handlers.GetCategoryById))
	http.Handle("POST /admin/categories", admin(handlers.CreateCategory))
	http.Handle("PUT /admin/category/{id}", admin(handlers.UpdateCategory))
	http.Handle("DELETE /admin/category/{id}", admin(handlers.DeleteCategory))

	http.Handle("GET /admin/announcements/stats", admin(handlers.GetAnnouncementStats))
	http.Handle("GET /admin/announcements", admin(handlers.GetAnnouncements))
	http.Handle("GET /admin/announcement/{id}", admin(handlers.GetAnnouncementById))
	http.Handle("POST /admin/announcements", admin(handlers.CreateAnnouncement))
	http.Handle("PUT /admin/announcement/{id}", admin(handlers.UpdateAnnouncement))
	http.Handle("DELETE /admin/announcement/{id}", admin(handlers.DeleteAnnouncement))

	http.Handle("GET /events", auth(handlers.GetPublicEventsForUser))
	http.Handle("GET /user/events", auth(handlers.GetUserEvents))
	http.Handle("DELETE /user/event/{id}", perm("create_event", handlers.DeleteMyEvent))
	http.Handle("POST /user/event/{id}/register", perm("register_event", handlers.RegisterForEvent))
	http.Handle("DELETE /user/event/{id}/unregister", perm("register_event", handlers.UnregisterFromEvent))

	http.Handle("POST /reports", auth(handlers.SubmitReport))

	http.Handle("GET /admin/reports", admin(handlers.GetAdminReports))
	http.Handle("GET /admin/reports/stats", admin(handlers.GetAdminReportStats))
	http.Handle("GET /admin/reports/users", admin(handlers.GetUsersHistorySummary))
	http.Handle("PATCH /admin/report/{id}", admin(handlers.UpdateAdminReport))
	http.Handle("DELETE /admin/report/{id}/content", admin(handlers.SoftDeleteReportContent))
	http.Handle("POST /admin/user/{id}/sanctions", admin(handlers.CreateSanction))
	http.Handle("GET /admin/user/{id}/history", admin(handlers.GetUserHistory))

	http.Handle("GET /admin/events", admin(handlers.GetEvents))
	http.Handle("GET /admin/event/{id}", admin(handlers.GetEventById))
	http.Handle("POST /admin/events", admin(handlers.CreateEvent))
	http.Handle("PUT /admin/event/{id}", admin(handlers.UpdateEvent))
	http.Handle("DELETE /admin/event/{id}", admin(handlers.DeleteEvent))
}
