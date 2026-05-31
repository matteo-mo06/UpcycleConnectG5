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

func GetAllUsersWithRoles(search, status, role string, limit, offset int) ([]models.UserListItem, int, error) {
	where := `WHERE 1=1`
	var args []any

	if search != "" {
		where += " AND (u.first_name LIKE ? OR u.last_name LIKE ? OR u.email LIKE ? OR CONCAT(u.first_name, ' ', u.last_name) LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s, s, s)
	}
	if status != "" {
		where += " AND u.status = ?"
		args = append(args, status)
	}
	if role != "" {
		where += " AND EXISTS (SELECT 1 FROM USER_ROLE ur2 JOIN ROLE r2 ON r2.id_role = ur2.id_role WHERE ur2.id_user = u.id_user AND r2.name_role = ?)"
		args = append(args, role)
	}

	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM USER u "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT
			u.id_user, u.email, u.first_name, u.last_name,
			u.upcycling_score, u.premium, u.status, u.created_at,
			COALESCE(GROUP_CONCAT(r.name_role ORDER BY r.name_role SEPARATOR ','), '') AS roles
		FROM USER u
		LEFT JOIN USER_ROLE ur ON ur.id_user = u.id_user
		LEFT JOIN ROLE r ON r.id_role = ur.id_role
		`+where+` GROUP BY u.id_user ORDER BY u.created_at DESC LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		if rolesStr != "" {
			item.Roles = strings.Split(rolesStr, ",")
		} else {
			item.Roles = []string{}
		}
		users = append(users, item)
	}
	return users, total, nil
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
	// 1. Récupère les annonces du vendeur avant toute suppression
	rows, err := config.Conn.Query("SELECT id_announcement FROM USER_ANNOUNCEMENT WHERE id_user = ?", id)
	if err != nil {
		return err
	}
	var announcementIDs []string
	for rows.Next() {
		var aid string
		if err := rows.Scan(&aid); err != nil {
			rows.Close()
			return err
		}
		announcementIDs = append(announcementIDs, aid)
	}
	rows.Close()

	// 2. Nullifie les FK nullable qui pointent vers cet utilisateur
	nullifies := []string{
		"UPDATE EVENT        SET id_creator       = NULL WHERE id_creator       = ?",
		"UPDATE PROJECT      SET id_creator       = NULL WHERE id_creator       = ?",
		"UPDATE FORMATION    SET id_creator       = NULL WHERE id_creator       = ?",
		"UPDATE FORMATION    SET id_formateur     = NULL WHERE id_formateur     = ?",
		"UPDATE REPORT       SET id_reported_user = NULL WHERE id_reported_user = ?",
		"UPDATE REPORT       SET resolved_by      = NULL WHERE resolved_by      = ?",
		"UPDATE ANNOUNCEMENT SET id_buyer         = NULL WHERE id_buyer         = ?",
	}
	for _, q := range nullifies {
		if _, err := config.Conn.Exec(q, id); err != nil {
			return err
		}
	}

	// 3. Libère les casiers et supprime les annonces du vendeur
	for _, aid := range announcementIDs {
		for _, q := range []string{
			"UPDATE REPORT SET id_announcement = NULL WHERE id_announcement = ?",
			"DELETE FROM ANNOUNCEMENT_LOCKER WHERE id_announcement = ?",
		} {
			if _, err := config.Conn.Exec(q, aid); err != nil {
				return err
			}
		}
	}
	if _, err := config.Conn.Exec("DELETE FROM USER_ANNOUNCEMENT WHERE id_user = ?", id); err != nil {
		return err
	}
	for _, aid := range announcementIDs {
		if _, err := config.Conn.Exec("DELETE FROM ANNOUNCEMENT WHERE id_announcement = ?", aid); err != nil {
			return err
		}
	}

	// 4. Supprime les sanctions (doit précéder la suppression des REPORT)
	for _, q := range []string{
		"DELETE FROM USER_SANCTIONS WHERE id_user  = ?",
		"DELETE FROM USER_SANCTIONS WHERE id_admin = ?",
	} {
		if _, err := config.Conn.Exec(q, id); err != nil {
			return err
		}
	}

	// 5. Supprime les autres tables liées à l'utilisateur
	tables := []string{
		"PROFESSIONAL_REQUEST",
		"DOCUMENT",
		"REPORT",
		"USER_ROLE",
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

	_, err = config.Conn.Exec("DELETE FROM USER WHERE id_user = ?", id)
	return err
}

func GetUserStats(userID string) (models.UserStats, error) {
	var s models.UserStats

	err := config.Conn.QueryRow("SELECT upcycling_score FROM USER WHERE id_user = ?", userID).Scan(&s.UpcyclingScore)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow(`
		SELECT COUNT(*) FROM ANNOUNCEMENT a
		JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		WHERE ua.id_user = ? AND a.state_annoucement = 'Active' AND a.request = 0`, userID).Scan(&s.ActiveAnnouncements)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow(`
		SELECT COUNT(*) FROM ANNOUNCEMENT a
		JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		LEFT JOIN ANNOUNCEMENT_LOCKER al ON al.id_announcement = a.id_announcement
		WHERE ua.id_user = ? AND a.request = 1 AND al.id_locker IS NULL`, userID).Scan(&s.PendingDeposits)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow(`
		SELECT COUNT(*) FROM USER_EVENT_INSCRIPTION uei
		JOIN EVENT e ON e.id_event = uei.id_event
		WHERE uei.id_user = ? AND e.date_event > NOW()`, userID).Scan(&s.UpcomingEvents)
	return s, err
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
