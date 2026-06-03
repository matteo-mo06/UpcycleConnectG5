package db

import (
	"errors"

	"github.com/google/uuid"
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
		err := rows.Scan(&r.Id_role, &r.Name_role)
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func GetRoleByName(name string) (models.Role, error) {
	var r models.Role
	err := config.Conn.QueryRow("SELECT id_role, name_role FROM ROLE WHERE name_role = ?", name).Scan(&r.Id_role, &r.Name_role)
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
		err := rows.Scan(&r.Id_role, &r.Name_role)
		if err != nil {
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
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

func AddUserRole(userID, roleID string) error {
	_, err := config.Conn.Exec(`
		INSERT IGNORE INTO USER_ROLE (
			id_user,
			id_role
		)
		VALUES (?, ?)
	`,
		userID,
		roleID,
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

func CreateRole(name string) (models.Role, error) {
	r := models.Role{
		Id_role:   uuid.New().String(),
		Name_role: name,
	}
	_, err := config.Conn.Exec(
		"INSERT INTO ROLE (id_role, name_role) VALUES (?, ?)",
		r.Id_role, r.Name_role,
	)
	return r, err
}

func UpdateRole(id, name string) error {
	_, err := config.Conn.Exec(
		"UPDATE ROLE SET name_role = ? WHERE id_role = ?",
		name, id,
	)
	return err
}

var protectedRoles = map[string]bool{"admin": true, "artisan": true, "salarie": true, "user": true}

func DeleteRole(id string) error {
	var name string
	if err := config.Conn.QueryRow("SELECT name_role FROM ROLE WHERE id_role = ?", id).Scan(&name); err != nil {
		return err
	}
	if protectedRoles[name] {
		return errors.New("ce rôle est protégé et ne peut pas être supprimé")
	}

	var count int
	if err := config.Conn.QueryRow(
		"SELECT COUNT(*) FROM USER_ROLE WHERE id_role = ?", id,
	).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("ce rôle est assigné à des utilisateurs")
	}
	if err := config.Conn.QueryRow(
		"SELECT COUNT(*) FROM ROLE_PERMISSION WHERE role_id = ?", id,
	).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("ce rôle possède des permissions, retirez-les d'abord")
	}
	_, err := config.Conn.Exec("DELETE FROM ROLE WHERE id_role = ?", id)
	return err
}
