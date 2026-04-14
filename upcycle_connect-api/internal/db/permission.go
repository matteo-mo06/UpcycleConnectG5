package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

// GetAllPermissions retourne toutes les permissions — utilisé par l'admin
func GetAllPermissions() ([]models.Permission, error) {
	rows, err := config.Conn.Query(`
		SELECT id, name, domain, created_at
		FROM PERMISSION
		ORDER BY domain, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.Id, &p.Name, &p.Domain, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

// GetPermissionsByRole retourne les permissions d'un rôle donné
func GetPermissionsByRole(roleID string) ([]models.Permission, error) {
	rows, err := config.Conn.Query(`
		SELECT p.id, p.name, p.domain, p.created_at
		FROM PERMISSION p
		INNER JOIN ROLE_PERMISSION rp ON rp.permission_id = p.id
		WHERE rp.role_id = ?
		ORDER BY p.domain, p.name`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.Id, &p.Name, &p.Domain, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

// GetUserPermissions retourne les noms de toutes les permissions d'un user (via ses rôles)
func GetUserPermissions(userID string) ([]string, error) {
	rows, err := config.Conn.Query(`
		SELECT DISTINCT p.name
		FROM PERMISSION p
		INNER JOIN ROLE_PERMISSION rp ON rp.permission_id = p.id
		INNER JOIN USER_ROLE ur       ON ur.id_role = rp.role_id
		WHERE ur.id_user = ?
		ORDER BY p.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}
	return perms, nil
}

// AddRolePermission assigne une permission à un rôle
func AddRolePermission(roleID, permissionID string) error {
	_, err := config.Conn.Exec(`
		INSERT IGNORE INTO ROLE_PERMISSION (id, role_id, permission_id)
		VALUES (UUID(), ?, ?)`,
		roleID, permissionID,
	)
	return err
}

// RemoveRolePermission retire une permission d'un rôle
func RemoveRolePermission(roleID, permissionID string) error {
	_, err := config.Conn.Exec(`
		DELETE FROM ROLE_PERMISSION
		WHERE role_id = ? AND permission_id = ?`,
		roleID, permissionID,
	)
	return err
}
