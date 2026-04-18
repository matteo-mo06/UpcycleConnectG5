package db

import (
	"strings"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetUserById(id string) (models.User, error) {
	var user models.User
	err := config.Conn.QueryRow(`
		SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium, status, created_at
		FROM USER
		WHERE id_user = ?`, id,
	).Scan(
		&user.Id_user,
		&user.Email,
		&user.Password_user,
		&user.First_name,
		&user.Last_name,
		&user.Upcycling_score,
		&user.Premium,
		&user.Status,
		&user.Created_at,
	)
	return user, err
}

func GetAllUsersByNameOrFirstName(firstName string, lastName string) ([]models.User, error) {
	query := `
		SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium, status, created_at
		FROM USER
		WHERE 1=1`
	var args []any

	if firstName != "" {
		query += " AND first_name LIKE ?"
		args = append(args, "%"+firstName+"%")
	}
	if lastName != "" {
		query += " AND last_name LIKE ?"
		args = append(args, "%"+lastName+"%")
	}
	query += " ORDER BY created_at DESC"

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Id_user,
			&user.Email,
			&user.Password_user,
			&user.First_name,
			&user.Last_name,
			&user.Upcycling_score,
			&user.Premium,
			&user.Status,
			&user.Created_at,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetAllUsersWithRoles(firstName string, lastName string) ([]models.UserListItem, error) {
	query := `
		SELECT
			u.id_user, u.email, u.first_name, u.last_name,
			u.upcycling_score, u.premium, u.status, u.created_at,
			COALESCE(GROUP_CONCAT(r.name_role ORDER BY r.name_role SEPARATOR ','), '') AS roles
		FROM USER u
		LEFT JOIN USER_ROLE ur ON ur.id_user = u.id_user
		LEFT JOIN ROLE r ON r.id_role = ur.id_role
		WHERE 1=1`
	var args []any

	if firstName != "" {
		query += " AND u.first_name LIKE ?"
		args = append(args, "%"+firstName+"%")
	}
	if lastName != "" {
		query += " AND u.last_name LIKE ?"
		args = append(args, "%"+lastName+"%")
	}
	query += " GROUP BY u.id_user ORDER BY u.created_at DESC"

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserListItem
	for rows.Next() {
		var item models.UserListItem
		var rolesStr string
		err := rows.Scan(
			&item.Id_user,
			&item.Email,
			&item.First_name,
			&item.Last_name,
			&item.Upcycling_score,
			&item.Premium,
			&item.Status,
			&item.Created_at,
			&rolesStr,
		)
		if err != nil {
			return nil, err
		}
		if rolesStr != "" {
			item.Roles = strings.Split(rolesStr, ",")
		} else {
			item.Roles = []string{}
		}
		users = append(users, item)
	}
	return users, nil
}

func CreateUser(user models.User) error {
	_, err := config.Conn.Exec(`
		INSERT INTO USER (id_user, email, password_user, first_name, last_name, upcycling_score, premium)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Id_user,
		user.Email,
		user.Password_user,
		user.First_name,
		user.Last_name,
		user.Upcycling_score,
		user.Premium,
	)
	return err
}

func UpdateUser(user models.User) error {
	_, err := config.Conn.Exec(`
		UPDATE USER SET
			email = ?,
			first_name = ?,
			last_name = ?,
			upcycling_score = ?,
			premium = ?
		WHERE id_user = ?`,
		user.Email,
		user.First_name,
		user.Last_name,
		user.Upcycling_score,
		user.Premium,
		user.Id_user,
	)
	return err
}

func UpdateUserStatus(id string, status string) error {
	_, err := config.Conn.Exec(`
		UPDATE USER SET status = ? WHERE id_user = ?`,
		status, id,
	)
	return err
}

func DeleteUser(id string) error {
	nullifies := []string{
		"UPDATE EVENT     SET id_creator   = NULL WHERE id_creator   = ?",
		"UPDATE PROJECT   SET id_creator   = NULL WHERE id_creator   = ?",
		"UPDATE FORMATION SET id_creator   = NULL WHERE id_creator   = ?",
		"UPDATE FORMATION SET id_formateur = NULL WHERE id_formateur = ?",
	}
	for _, q := range nullifies {
		if _, err := config.Conn.Exec(q, id); err != nil {
			return err
		}
	}

	tables := []string{
		"PROFESSIONAL_REQUEST",
		"DOCUMENT",
		"REPORT",
		"USER_ROLE",
		"USER_ANNOUNCEMENT",
		"USER_ADVICE",
		"USER_TOPIC",
		"USER_PROJECT_INSCRIPTION",
		"USER_PROJECT_UPDOWN",
		"USER_FORMATION_INSCRIPTION",
		"USER_FORMATION_INSCRIPTION_PAYEMENT",
		"USER_EVENT_INSCRIPTION",
		"USER_TICKET",
		"USER_NOTIFICATION",
		"USER_PAYEMENT",
		"USER_SUBSCRIPTION",
	}

	for _, table := range tables {
		if _, err := config.Conn.Exec("DELETE FROM "+table+" WHERE id_user = ?", id); err != nil {
			return err
		}
	}

	_, err := config.Conn.Exec("DELETE FROM USER WHERE id_user = ?", id)
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := config.Conn.QueryRow(`
		SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium, status, created_at
		FROM USER
		WHERE email = ?`, email,
	).Scan(
		&user.Id_user,
		&user.Email,
		&user.Password_user,
		&user.First_name,
		&user.Last_name,
		&user.Upcycling_score,
		&user.Premium,
		&user.Status,
		&user.Created_at,
	)
	return user, err
}
