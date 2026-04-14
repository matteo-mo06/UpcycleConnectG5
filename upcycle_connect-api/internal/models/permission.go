package models

type Permission struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at"`
}

type RolePermission struct {
	Id           string `json:"id"`
	RoleId       string `json:"role_id"`
	PermissionId string `json:"permission_id"`
}
