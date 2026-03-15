package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllRoles() ([]models.Role, error) {
	rows, err := config.Conn.Query("SELECT id_role, name_role FROM ROLE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.Id_role, &r.Name_role); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func GetRoleByName(name string) (models.Role, error) {
	var r models.Role
	row := config.Conn.QueryRow("SELECT id_role, name_role FROM ROLE WHERE name_role = ?", name)
	err := row.Scan(&r.Id_role, &r.Name_role)
	return r, err
}

func GetUserRoles(userID string) ([]models.Role, error) {
	rows, err := config.Conn.Query(`
		SELECT r.id_role, r.name_role
		FROM ROLE r
		JOIN USER_ROLE ur ON r.id_role = ur.id_role
		WHERE ur.id_user = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.Id_role, &r.Name_role); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func GetUserRoleNames(userID string) ([]string, error) {
	rows, err := config.Conn.Query(`
		SELECT r.name_role
		FROM ROLE r
		JOIN USER_ROLE ur ON r.id_role = ur.id_role
		WHERE ur.id_user = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

func AddUserRole(userID, roleID string) error {
	_, err := config.Conn.Exec(
		"INSERT IGNORE INTO USER_ROLE (id_user, id_role) VALUES (?, ?)",
		userID, roleID,
	)
	return err
}

func RemoveUserRole(userID, roleID string) (bool, error) {
	result, err := config.Conn.Exec(
		"DELETE FROM USER_ROLE WHERE id_user = ? AND id_role = ?",
		userID, roleID,
	)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	return rows > 0, err
}
