package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetFormations(userID, search, level, idCategory string, limit, offset int) ([]models.Formation, int, error) {
	where := `WHERE f.deleted_at IS NULL AND (f.status = 'approved' OR f.id_creator = ?)`
	args := []any{userID}

	if search != "" {
		where += " AND f.title_formation LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if level != "" {
		where += " AND f.level = ?"
		args = append(args, level)
	}
	if idCategory != "" {
		where += " AND f.id_category = ?"
		args = append(args, idCategory)
	}

	countArgs := make([]any, len(args))
	copy(countArgs, args)
	var total int
	if err := config.Conn.QueryRow("SELECT COUNT(DISTINCT f.id_formation) FROM formation f "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	isRegArgs := append([]any{userID}, args...)
	isRegArgs = append(isRegArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT f.id_formation, f.title_formation, f.description_formation, f.date_formation,
		       f.location_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       COUNT(ins.id_user) AS inscription_count,
		       EXISTS(SELECT 1 FROM user_formation_inscription r WHERE r.id_user = ? AND r.id_formation = f.id_formation) AS is_registered,
		       f.created_at
		FROM formation f
		LEFT JOIN category c ON c.id_category = f.id_category
		LEFT JOIN user u ON u.id_user = f.id_creator
		LEFT JOIN user fmt ON fmt.id_user = f.id_formateur
		LEFT JOIN user_formation_inscription ins ON ins.id_formation = f.id_formation
		`+where+`
		GROUP BY f.id_formation
		ORDER BY f.date_formation ASC
		LIMIT ? OFFSET ?`, isRegArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Formation
	for rows.Next() {
		var f models.Formation
		err := rows.Scan(
			&f.IdFormation, &f.TitleFormation, &f.DescriptionFormation, &f.DateFormation,
			&f.LocationFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.InscriptionCount, &f.IsRegistered,
			&f.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, f)
	}
	return list, total, nil
}

func GetAllFormations(search, status string, limit, offset int) ([]models.Formation, int, error) {
	where := "WHERE f.deleted_at IS NULL"
	var args []any

	if search != "" {
		where += " AND f.title_formation LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if status != "" {
		where += " AND f.status = ?"
		args = append(args, status)
	}

	countArgs := make([]any, len(args))
	copy(countArgs, args)

	var total int
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM formation f "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT f.id_formation, f.title_formation, f.description_formation, f.date_formation,
		       f.location_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       COUNT(ins.id_user) AS inscription_count,
		       f.created_at
		FROM formation f
		LEFT JOIN category c ON c.id_category = f.id_category
		LEFT JOIN user u ON u.id_user = f.id_creator
		LEFT JOIN user fmt ON fmt.id_user = f.id_formateur
		LEFT JOIN user_formation_inscription ins ON ins.id_formation = f.id_formation
		`+where+`
		GROUP BY f.id_formation
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Formation
	for rows.Next() {
		var f models.Formation
		err := rows.Scan(
			&f.IdFormation, &f.TitleFormation, &f.DescriptionFormation, &f.DateFormation,
			&f.LocationFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.InscriptionCount,
			&f.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, f)
	}
	return list, total, nil
}

func GetFormationById(id string) (models.Formation, error) {
	var f models.Formation
	err := config.Conn.QueryRow(`
		SELECT f.id_formation, f.title_formation, f.description_formation, f.date_formation,
		       f.location_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       f.created_at
		FROM formation f
		LEFT JOIN category c ON c.id_category = f.id_category
		LEFT JOIN user u ON u.id_user = f.id_creator
		LEFT JOIN user fmt ON fmt.id_user = f.id_formateur
		WHERE f.id_formation = ? AND f.deleted_at IS NULL`, id,
	).Scan(
		&f.IdFormation, &f.TitleFormation, &f.DescriptionFormation, &f.DateFormation,
		&f.LocationFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
		&f.RejectionReason, &f.IdCategory, &f.CategoryName,
		&f.IdCreator, &f.CreatorName,
		&f.IdFormateur, &f.FormateurName,
		&f.CreatedAt,
	)
	return f, err
}

func CreateFormation(f models.Formation) error {
	_, err := config.Conn.Exec(`
		INSERT INTO formation (
			id_formation, title_formation, description_formation, date_formation,
			location_formation, capacity, level, duration_hours, status,
			id_category, id_creator, id_formateur, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pending', ?, ?, ?, NOW())`,
		f.IdFormation, f.TitleFormation, f.DescriptionFormation, f.DateFormation,
		f.LocationFormation, f.Capacity, f.Level, f.DurationHours,
		f.IdCategory, f.IdCreator, f.IdFormateur,
	)
	return err
}

func UpdateFormation(f models.Formation) error {
	_, err := config.Conn.Exec(`
		UPDATE formation SET
			title_formation = ?,
			description_formation = ?,
			date_formation = ?,
			location_formation = ?,
			capacity = ?,
			level = ?,
			duration_hours = ?,
			id_category = ?,
			status = 'pending',
			rejection_reason = NULL
		WHERE id_formation = ?`,
		f.TitleFormation, f.DescriptionFormation, f.DateFormation,
		f.LocationFormation, f.Capacity, f.Level, f.DurationHours,
		f.IdCategory, f.IdFormation,
	)
	return err
}

func UpdateFormationAdmin(f models.Formation) error {
	_, err := config.Conn.Exec(`
		UPDATE formation SET
			title_formation = ?,
			description_formation = ?,
			date_formation = ?,
			location_formation = ?,
			capacity = ?,
			level = ?,
			duration_hours = ?,
			id_category = ?,
			status = ?
		WHERE id_formation = ?`,
		f.TitleFormation, f.DescriptionFormation, f.DateFormation,
		f.LocationFormation, f.Capacity, f.Level, f.DurationHours,
		f.IdCategory, f.Status, f.IdFormation,
	)
	return err
}

func ApproveFormation(id string) error {
	_, err := config.Conn.Exec(
		"UPDATE formation SET status = 'approved', rejection_reason = NULL WHERE id_formation = ? AND deleted_at IS NULL", id,
	)
	return err
}

func RejectFormation(id string, reason *string) error {
	_, err := config.Conn.Exec(
		"UPDATE formation SET status = 'rejected', rejection_reason = ? WHERE id_formation = ? AND deleted_at IS NULL", reason, id,
	)
	return err
}

func DeleteFormation(id string) error {
	_, err := config.Conn.Exec("UPDATE formation SET deleted_at = NOW() WHERE id_formation = ?", id)
	return err
}

func RegisterUserForFormation(userID, formationID string) error {
	_, err := config.Conn.Exec(
		"INSERT IGNORE INTO user_formation_inscription (id_user, id_formation) VALUES (?, ?)",
		userID, formationID,
	)
	return err
}

func UnregisterUserFromFormation(userID, formationID string) error {
	_, err := config.Conn.Exec(
		"DELETE FROM user_formation_inscription WHERE id_user = ? AND id_formation = ?",
		userID, formationID,
	)
	return err
}

func GetUserRegisteredFormations(userID string) ([]models.Formation, error) {
	rows, err := config.Conn.Query(`
		SELECT f.id_formation, f.title_formation, f.description_formation, f.date_formation,
		       f.location_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       COUNT(ins.id_user) AS inscription_count,
		       f.created_at
		FROM formation f
		LEFT JOIN category c ON c.id_category = f.id_category
		LEFT JOIN user u ON u.id_user = f.id_creator
		LEFT JOIN user fmt ON fmt.id_user = f.id_formateur
		LEFT JOIN user_formation_inscription ins ON ins.id_formation = f.id_formation
		WHERE f.deleted_at IS NULL AND EXISTS (
			SELECT 1 FROM user_formation_inscription WHERE id_user = ? AND id_formation = f.id_formation
		)
		GROUP BY f.id_formation
		ORDER BY f.date_formation ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Formation
	for rows.Next() {
		var f models.Formation
		f.IsRegistered = true
		err := rows.Scan(
			&f.IdFormation, &f.TitleFormation, &f.DescriptionFormation, &f.DateFormation,
			&f.LocationFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.InscriptionCount,
			&f.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, nil
}

func GetMyCreatedFormations(creatorID string) ([]models.Formation, error) {
	rows, err := config.Conn.Query(`
		SELECT f.id_formation, f.title_formation, f.description_formation, f.date_formation,
		       f.location_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       COUNT(ins.id_user) AS inscription_count,
		       f.created_at
		FROM formation f
		LEFT JOIN category c ON c.id_category = f.id_category
		LEFT JOIN user u ON u.id_user = f.id_creator
		LEFT JOIN user fmt ON fmt.id_user = f.id_formateur
		LEFT JOIN user_formation_inscription ins ON ins.id_formation = f.id_formation
		WHERE f.id_creator = ? AND f.deleted_at IS NULL
		GROUP BY f.id_formation
		ORDER BY f.created_at DESC`, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Formation
	for rows.Next() {
		var f models.Formation
		err := rows.Scan(
			&f.IdFormation, &f.TitleFormation, &f.DescriptionFormation, &f.DateFormation,
			&f.LocationFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.InscriptionCount,
			&f.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, nil
}
