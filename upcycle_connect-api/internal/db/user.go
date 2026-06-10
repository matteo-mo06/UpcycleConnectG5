package db

import (
	"strings"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetUserById(id string) (models.User, error) {
	var user models.User
	err := config.Conn.QueryRow(`
		SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium, status, created_at, avatar_url, tutorial_done
		FROM USER
		WHERE id_user = ? AND deleted_at IS NULL`, id,
	).Scan(
		&user.IdUser,
		&user.Email,
		&user.PasswordUser,
		&user.FirstName,
		&user.LastName,
		&user.UpcyclingScore,
		&user.Premium,
		&user.Status,
		&user.CreatedAt,
		&user.AvatarUrl,
		&user.TutorialDone,
	)
	return user, err
}

func GetAllUsersByNameOrFirstName(firstName string, lastName string) ([]models.User, error) {
	query := `
		SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium, status, created_at
		FROM USER
		WHERE deleted_at IS NULL`
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
			&user.IdUser,
			&user.Email,
			&user.PasswordUser,
			&user.FirstName,
			&user.LastName,
			&user.UpcyclingScore,
			&user.Premium,
			&user.Status,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetAllUsersWithRoles(search, status, role string, limit, offset int) ([]models.UserListItem, int, error) {
	where := `WHERE u.deleted_at IS NULL`
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
			&item.IdUser,
			&item.Email,
			&item.FirstName,
			&item.LastName,
			&item.UpcyclingScore,
			&item.Premium,
			&item.Status,
			&item.CreatedAt,
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
		user.IdUser,
		user.Email,
		user.PasswordUser,
		user.FirstName,
		user.LastName,
		user.UpcyclingScore,
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
		user.FirstName,
		user.LastName,
		user.UpcyclingScore,
		user.Premium,
		user.IdUser,
	)
	return err
}

func UpdateUserAvatar(id, avatarURL string) error {
	_, err := config.Conn.Exec(
		`UPDATE USER SET avatar_url = ? WHERE id_user = ? AND deleted_at IS NULL`,
		avatarURL, id)
	return err
}

func UpdateUserPassword(id, hashedPassword string) error {
	_, err := config.Conn.Exec(
		`UPDATE USER SET password_user = ? WHERE id_user = ? AND deleted_at IS NULL`,
		hashedPassword, id)
	return err
}

func UpdateUserProfile(id, firstName, lastName, email string) error {
	_, err := config.Conn.Exec(`
		UPDATE USER SET first_name = ?, last_name = ?, email = ?
		WHERE id_user = ? AND deleted_at IS NULL`,
		firstName, lastName, email, id)
	return err
}

func UpdateUserStatus(id string, status string) error {
	_, err := config.Conn.Exec(`
		UPDATE USER SET status = ? WHERE id_user = ?`,
		status, id,
	)
	return err
}

func MarkTutorialDone(id string) error {
	_, err := config.Conn.Exec(
		`UPDATE USER SET tutorial_done = 1 WHERE id_user = ? AND deleted_at IS NULL`, id)
	return err
}

func DeleteUser(id string) error {
	_, err := config.Conn.Exec(`
		UPDATE USER SET
			email          = CONCAT('deleted_', id_user, '@deleted.invalid'),
			password_user  = '',
			first_name     = 'Utilisateur',
			last_name      = 'supprimé',
			deleted_at     = NOW()
		WHERE id_user = ?`, id)
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
		SELECT id_user, email, password_user, first_name, last_name, upcycling_score, premium, status, created_at, avatar_url, tutorial_done
		FROM USER
		WHERE email = ? AND deleted_at IS NULL`, email,
	).Scan(
		&user.IdUser,
		&user.Email,
		&user.PasswordUser,
		&user.FirstName,
		&user.LastName,
		&user.UpcyclingScore,
		&user.Premium,
		&user.Status,
		&user.CreatedAt,
		&user.AvatarUrl,
		&user.TutorialDone,
	)
	return user, err
}


func SaveOnesignalPlayerID(userID, playerID string) error {
	_, err := config.Conn.Exec(
		"UPDATE user SET onesignal_player_id = ? WHERE id_user = ?",
		playerID, userID,
	)
	return err
}

func GetOnesignalPlayerID(userID string) string {
	var pid *string
	_ = config.Conn.QueryRow(
		"SELECT onesignal_player_id FROM user WHERE id_user = ?", userID,
	).Scan(&pid)
	if pid == nil {
		return ""
	}
	return *pid
}

func GetEventCreatorID(eventID string) string {
	var id string
	_ = config.Conn.QueryRow(
		"SELECT id_creator FROM event WHERE id_event = ?", eventID,
	).Scan(&id)
	return id
}

func GetFormationCreatorID(formationID string) string {
	var id string
	_ = config.Conn.QueryRow(
		"SELECT id_creator FROM formation WHERE id_formation = ?", formationID,
	).Scan(&id)
	return id
}
