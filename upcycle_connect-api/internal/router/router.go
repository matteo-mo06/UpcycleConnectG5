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

	http.Handle("GET /categories", auth(handlers.GetCategories))

	http.HandleFunc("POST /auth/login", handlers.Login)
	http.HandleFunc("POST /auth/register", handlers.CreateUser)

	http.Handle("PATCH /user/profile", auth(handlers.UpdateMyProfile))
	http.Handle("PATCH /user/password", auth(handlers.UpdateMyPassword))
	http.Handle("PATCH /user/avatar", auth(handlers.UpdateMyAvatar))
	http.Handle("DELETE /user/account", auth(handlers.DeleteMyAccount))
	http.Handle("GET /user/me", auth(handlers.GetMe))
	http.Handle("GET /user/stats", auth(handlers.GetMyStats))
	http.Handle("GET /user/score-breakdown", auth(handlers.GetMyScoreBreakdown))
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
	http.Handle("GET /user/professional-request", auth(handlers.GetMyProfessionalRequest))

	http.Handle("GET /admin/stats", admin(handlers.GetStats))

	http.Handle("GET /admin/users", admin(handlers.GetUsers))
	http.Handle("GET /admin/user/{id}", admin(handlers.GetUserById))
	http.Handle("PUT /admin/user/{id}", admin(handlers.UpdateUser))
	http.Handle("PATCH /admin/user/{id}/status", admin(handlers.UpdateUserStatus))
	http.Handle("DELETE /admin/user/{id}", admin(handlers.DeleteUser))

	http.Handle("GET /admin/roles", perm("manage_roles", handlers.GetRoles))
	http.Handle("POST /admin/roles", perm("manage_roles", handlers.CreateRole))
	http.Handle("PUT /admin/role/{id}", perm("manage_roles", handlers.UpdateRole))
	http.Handle("DELETE /admin/role/{id}", perm("manage_roles", handlers.DeleteRole))
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
	http.Handle("PATCH /admin/announcement/{id}/approve", perm("manage_announcements", handlers.ApproveAnnouncement))
	http.Handle("PATCH /admin/announcement/{id}/reject", perm("manage_announcements", handlers.RejectAnnouncement))

	http.Handle("GET /events", auth(handlers.GetPublicEventsForUser))
	http.Handle("POST /events", perm("create_event", handlers.CreateEvent))
	http.Handle("GET /user/events", auth(handlers.GetUserEvents))
	http.Handle("DELETE /user/event/{id}", perm("create_event", handlers.DeleteMyEvent))
	http.Handle("POST /user/event/{id}/register", perm("register_event", handlers.RegisterForEvent))
	http.Handle("DELETE /user/event/{id}/unregister", perm("register_event", handlers.UnregisterFromEvent))

	http.Handle("GET /forum/topics", auth(handlers.GetForumTopics))
	http.Handle("GET /forum/topics/{id}", auth(handlers.GetForumTopic))
	http.Handle("POST /forum/topics", perm("create_topic", handlers.CreateForumTopic))
	http.Handle("PATCH /forum/topics/{id}", auth(handlers.UpdateForumTopic))
	http.Handle("DELETE /forum/topics/{id}", auth(handlers.DeleteForumTopic))
	http.Handle("POST /forum/topics/{id}/posts", perm("create_post", handlers.AddForumPost))
	http.Handle("PATCH /forum/topics/{id}/posts/{post_id}", auth(handlers.UpdateForumPost))
	http.Handle("DELETE /forum/topics/{id}/posts/{post_id}", auth(handlers.DeleteForumPost))

	http.Handle("POST /reports", auth(handlers.SubmitReport))

	http.Handle("GET /admin/reports", admin(handlers.GetAdminReports))
	http.Handle("GET /admin/reports/stats", admin(handlers.GetAdminReportStats))
	http.Handle("GET /admin/reports/users", admin(handlers.GetUsersHistorySummary))
	http.Handle("PATCH /admin/report/{id}", admin(handlers.UpdateAdminReport))
	http.Handle("DELETE /admin/report/{id}/content", admin(handlers.SoftDeleteReportContent))
	http.Handle("POST /admin/user/{id}/sanctions", admin(handlers.CreateSanction))
	http.Handle("GET /admin/user/{id}/history", admin(handlers.GetUserHistory))

	http.Handle("GET /admin/score-actions", admin(handlers.GetScoreActions))
	http.Handle("PUT /admin/score-action/{action_type}", admin(handlers.UpdateScoreAction))

	http.Handle("GET /admin/lockers", perm("manage_lockers", handlers.GetLockers))
	http.Handle("POST /admin/lockers", perm("manage_lockers", handlers.CreateLocker))
	http.Handle("DELETE /admin/locker/{id}", perm("manage_lockers", handlers.DeleteLocker))
	http.Handle("PATCH /admin/locker/{id}/access-code", perm("manage_lockers", handlers.UpdateLockerAccessCode))

	http.Handle("GET /user/my-events", perm("create_event", handlers.GetMyCreatedEvents))

	http.Handle("GET /admin/events", admin(handlers.GetEvents))
	http.Handle("GET /admin/event/{id}", admin(handlers.GetEventById))
	http.Handle("POST /admin/events", admin(handlers.CreateEventAdmin))
	http.Handle("PUT /admin/event/{id}", admin(handlers.UpdateEvent))
	http.Handle("DELETE /admin/event/{id}", admin(handlers.DeleteEvent))
	http.Handle("PATCH /admin/event/{id}/approve", admin(handlers.ApproveEvent))
	http.Handle("PATCH /admin/event/{id}/reject", admin(handlers.RejectEvent))

	http.Handle("GET /formations", auth(handlers.GetFormations))
	http.Handle("GET /formations/pending", perm("manage_formations", handlers.GetPendingFormations))
	http.Handle("GET /formations/{id}", auth(handlers.GetFormationById))
	http.Handle("POST /formations", perm("create_formation", handlers.CreateFormation))
	http.Handle("PATCH /formations/{id}", auth(handlers.UpdateMyFormation))
	http.Handle("DELETE /user/formation/{id}", auth(handlers.DeleteMyFormation))
	http.Handle("PATCH /formations/{id}/approve", perm("manage_formations", handlers.ApproveFormation))
	http.Handle("PATCH /formations/{id}/reject", perm("manage_formations", handlers.RejectFormation))
	http.Handle("POST /user/formation/{id}/register", perm("register_formation", handlers.RegisterForFormation))
	http.Handle("DELETE /user/formation/{id}/unregister", perm("register_formation", handlers.UnregisterFromFormation))
	http.Handle("GET /user/formations", auth(handlers.GetUserFormations))
	http.Handle("GET /user/my-formations", perm("create_formation", handlers.GetMyCreatedFormations))

	http.Handle("GET /admin/formations", admin(handlers.GetAllFormationsAdmin))
	http.Handle("GET /admin/formation/{id}", admin(handlers.GetFormationByIdAdmin))
	http.Handle("PUT /admin/formation/{id}", admin(handlers.UpdateFormationAdmin))
	http.Handle("DELETE /admin/formation/{id}", admin(handlers.DeleteFormationAdmin))

	http.Handle("GET /projects", auth(handlers.GetProjects))
	http.Handle("GET /projects/pending", perm("manage_projects", handlers.GetPendingProjects))
	http.Handle("GET /projects/{id}", auth(handlers.GetProjectById))
	http.Handle("POST /projects", perm("create_project", handlers.CreateProject))
	http.Handle("PATCH /projects/{id}", auth(handlers.UpdateMyProject))
	http.Handle("DELETE /user/project/{id}", perm("create_project", handlers.DeleteMyProject))
	http.Handle("PATCH /projects/{id}/approve", perm("manage_projects", handlers.ApproveProject))
	http.Handle("PATCH /projects/{id}/reject", perm("manage_projects", handlers.RejectProject))
	http.Handle("POST /projects/{id}/join", perm("register_project", handlers.RegisterForProject))
	http.Handle("DELETE /projects/{id}/join", perm("register_project", handlers.UnregisterFromProject))
	http.Handle("GET /user/projects", auth(handlers.GetUserProjects))
	http.Handle("GET /user/my-projects", perm("create_project", handlers.GetMyCreatedProjects))

	http.Handle("PUT /projects/{id}/vote", auth(handlers.VoteProject))
	http.Handle("DELETE /projects/{id}/vote", auth(handlers.RemoveProjectVote))
	http.Handle("DELETE /projects/{id}", auth(handlers.ModerateDeleteProject))

	http.Handle("GET /admin/projects/stats", admin(handlers.GetProjectStats))
	http.Handle("GET /admin/projects", admin(handlers.GetAllProjectsAdmin))
	http.Handle("GET /admin/project/{id}", admin(handlers.GetProjectByIdAdmin))
	http.Handle("PUT /admin/project/{id}", admin(handlers.UpdateProjectAdmin))
	http.Handle("DELETE /admin/project/{id}", admin(handlers.DeleteProjectAdmin))
	http.Handle("PATCH /admin/project/{id}/approve", admin(handlers.ApproveProject))
	http.Handle("PATCH /admin/project/{id}/reject", admin(handlers.RejectProject))
}
