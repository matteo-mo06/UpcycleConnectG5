package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetUserById(id string) (models.User, error) {
	var user models.User
	err := config.Conn.QueryRow(`
		SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium
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
	)
	return user, err
}

func GetAllUsersByNameOrFirstName(firstName string, lastName string) ([]models.User, error) {
	query := `SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium
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
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func CreateUser(user models.User) error {
	_, err := config.Conn.Exec(`
		INSERT INTO USER (
			id_user,
			email,
			password_user,
			first_name,
			last_name,
			upcycling_score,
			premium
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
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
		WHERE id_user = ?
	`,
		user.Email,
		user.First_name,
		user.Last_name,
		user.Upcycling_score,
		user.Premium,
		user.Id_user,
	)
	return err
}

func DeleteUser(id string) error {
	tables := []string{
		"PROFESSIONAL_REQUEST",
		"DOCUMENT",
		"REPORT_ANNOUCEMENT",
		"REPORT_TOPIC",
		"REPORT_POST",
		"USER_ROLE",
		"USER_ANNOUNCEMENT",
		"USER_ADVICE",
		"USER_TOPIC",
		"USER_PROJECT_CREATE",
		"USER_PROJECT_INSCRIPTION",
		"USER_PROJECT_UPDOWN",
		"USER_FORMATION_FORMATER",
		"USER_FORMATION_CREATE",
		"USER_FORMATION_INSCRIPTION",
		"USER_FORMATION_INSCRIPTION_PAYEMENT",
		"USER_EVENT",
		"USER_EVENT_INSCRIPTION",
		"USER_TICKET",
		"USER_NOTIFICATION",
		"USER_PAYEMENT",
		"USER_SUBSCRIPTION",
	}

	for _, table := range tables {
		_, err := config.Conn.Exec("DELETE FROM "+table+" WHERE id_user = ?", id)
		if err != nil {
			return err
		}
	}

	_, err := config.Conn.Exec("DELETE FROM USER WHERE id_user = ?", id)
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := config.Conn.QueryRow(`
		SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium
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
	)
	return user, err
}
