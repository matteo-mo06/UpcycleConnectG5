package db

import (
	"database/sql"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetFormationParticipants(formationID string) ([]models.Participant, error) {
	rows, err := config.Conn.Query(`
		SELECT u.id_user, u.first_name, u.last_name, u.avatar_url
		FROM user u
		JOIN user_formation_inscription fi ON fi.id_user = u.id_user
		WHERE fi.id_formation = ?`, formationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Participant
	for rows.Next() {
		var p models.Participant
		var avatarUrl sql.NullString
		if err := rows.Scan(&p.Id, &p.FirstName, &p.LastName, &avatarUrl); err != nil {
			return nil, err
		}
		if avatarUrl.Valid {
			p.AvatarUrl = &avatarUrl.String
		}
		list = append(list, p)
	}
	return list, nil
}

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
		       f.address_formation, f.city_formation, f.postal_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       f.syllabus, f.prerequisites, f.objectives,
		       COUNT(ins.id_user) AS inscription_count,
		       EXISTS(SELECT 1 FROM user_formation_inscription r WHERE r.id_user = ? AND r.id_formation = f.id_formation) AS is_registered,
		       f.created_at, f.price
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
			&f.AddressFormation, &f.CityFormation, &f.PostalFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.Syllabus, &f.Prerequisites, &f.Objectives,
			&f.InscriptionCount, &f.IsRegistered,
			&f.CreatedAt, &f.Price,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, f)
	}
	return list, total, nil
}

func GetAllFormations(search, status, level string, limit, offset int) ([]models.Formation, int, error) {
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
	if level != "" {
		where += " AND f.level = ?"
		args = append(args, level)
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
		       f.address_formation, f.city_formation, f.postal_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       f.syllabus, f.prerequisites, f.objectives,
		       COUNT(ins.id_user) AS inscription_count,
		       f.created_at, f.price
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
			&f.AddressFormation, &f.CityFormation, &f.PostalFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.Syllabus, &f.Prerequisites, &f.Objectives,
			&f.InscriptionCount,
			&f.CreatedAt, &f.Price,
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
		       f.address_formation, f.city_formation, f.postal_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       f.syllabus, f.prerequisites, f.objectives,
		       f.created_at, f.price
		FROM formation f
		LEFT JOIN category c ON c.id_category = f.id_category
		LEFT JOIN user u ON u.id_user = f.id_creator
		LEFT JOIN user fmt ON fmt.id_user = f.id_formateur
		WHERE f.id_formation = ? AND f.deleted_at IS NULL`, id,
	).Scan(
		&f.IdFormation, &f.TitleFormation, &f.DescriptionFormation, &f.DateFormation,
		&f.AddressFormation, &f.CityFormation, &f.PostalFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
		&f.RejectionReason, &f.IdCategory, &f.CategoryName,
		&f.IdCreator, &f.CreatorName,
		&f.IdFormateur, &f.FormateurName,
		&f.Syllabus, &f.Prerequisites, &f.Objectives,
		&f.CreatedAt, &f.Price,
	)
	return f, err
}

func GetFormationOwnerID(formationID string) (string, error) {
	var ownerID sql.NullString
	err := config.Conn.QueryRow(
		`SELECT id_creator FROM formation WHERE id_formation = ? AND deleted_at IS NULL`, formationID,
	).Scan(&ownerID)
	if err != nil {
		return "", err
	}
	return ownerID.String, nil
}

func CreateFormation(f models.Formation) error {
	status := f.Status
	if status == "" {
		status = "pending"
	}
	_, err := config.Conn.Exec(`
		INSERT INTO formation (
			id_formation, title_formation, description_formation, date_formation,
			address_formation, city_formation, postal_formation, capacity, level, duration_hours, status,
			id_category, id_creator, id_formateur, price, syllabus, prerequisites, objectives, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		f.IdFormation, f.TitleFormation, f.DescriptionFormation, f.DateFormation,
		f.AddressFormation, f.CityFormation, f.PostalFormation, f.Capacity, f.Level, f.DurationHours, status,
		f.IdCategory, f.IdCreator, f.IdFormateur, f.Price, f.Syllabus, f.Prerequisites, f.Objectives,
	)
	return err
}

func UpdateFormation(f models.Formation) error {
	_, err := config.Conn.Exec(`
		UPDATE formation SET
			title_formation = ?,
			description_formation = ?,
			date_formation = ?,
			address_formation = ?,
			city_formation = ?,
			postal_formation = ?,
			capacity = ?,
			level = ?,
			duration_hours = ?,
			id_category = ?,
			price = ?,
			syllabus = ?,
			prerequisites = ?,
			objectives = ?,
			status = 'pending',
			rejection_reason = NULL
		WHERE id_formation = ?`,
		f.TitleFormation, f.DescriptionFormation, f.DateFormation,
		f.AddressFormation, f.CityFormation, f.PostalFormation, f.Capacity, f.Level, f.DurationHours,
		f.IdCategory, f.Price, f.Syllabus, f.Prerequisites, f.Objectives, f.IdFormation,
	)
	return err
}

func UpdateFormationAdmin(f models.Formation) error {
	_, err := config.Conn.Exec(`
		UPDATE formation SET
			title_formation = ?,
			description_formation = ?,
			date_formation = ?,
			address_formation = ?,
			city_formation = ?,
			postal_formation = ?,
			capacity = ?,
			level = ?,
			duration_hours = ?,
			id_category = ?,
			price = ?,
			syllabus = ?,
			prerequisites = ?,
			objectives = ?,
			status = ?
		WHERE id_formation = ?`,
		f.TitleFormation, f.DescriptionFormation, f.DateFormation,
		f.AddressFormation, f.CityFormation, f.PostalFormation, f.Capacity, f.Level, f.DurationHours,
		f.IdCategory, f.Price, f.Syllabus, f.Prerequisites, f.Objectives, f.Status, f.IdFormation,
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

func GetDeletedFormations() ([]models.DeletedFormation, error) {
	rows, err := config.Conn.Query(`
		SELECT f.id_formation, f.title_formation, f.description_formation, f.date_formation,
		       f.address_formation, f.city_formation, f.postal_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category, c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       f.syllabus, f.prerequisites, f.objectives,
		       f.created_at, f.price, f.deleted_at
		FROM formation f
		LEFT JOIN category c ON c.id_category = f.id_category
		LEFT JOIN user u ON u.id_user = f.id_creator
		LEFT JOIN user fmt ON fmt.id_user = f.id_formateur
		WHERE f.deleted_at IS NOT NULL
		ORDER BY f.deleted_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DeletedFormation
	for rows.Next() {
		var f models.DeletedFormation
		if err := rows.Scan(
			&f.IdFormation, &f.TitleFormation, &f.DescriptionFormation, &f.DateFormation,
			&f.AddressFormation, &f.CityFormation, &f.PostalFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.Syllabus, &f.Prerequisites, &f.Objectives,
			&f.CreatedAt, &f.Price, &f.DeletedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, f)
	}

	for i := range list {
		steps, err := GetFormationSteps(list[i].IdFormation)
		if err != nil {
			return nil, err
		}
		if steps == nil {
			steps = []models.FormationStep{}
		}
		list[i].Steps = steps

		docs, err := GetDocumentsByCategory(list[i].IdFormation)
		if err != nil {
			return nil, err
		}
		if docs == nil {
			docs = []models.Document{}
		}
		list[i].Documents = docs
	}
	return list, nil
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
		       f.address_formation, f.city_formation, f.postal_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       COUNT(ins.id_user) AS inscription_count,
		       f.created_at, f.price
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
			&f.AddressFormation, &f.CityFormation, &f.PostalFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.InscriptionCount,
			&f.CreatedAt, &f.Price,
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
		       f.address_formation, f.city_formation, f.postal_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       COUNT(ins.id_user) AS inscription_count,
		       f.created_at, f.price
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
			&f.AddressFormation, &f.CityFormation, &f.PostalFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.InscriptionCount,
			&f.CreatedAt, &f.Price,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, nil
}

func GetMyCreatedFormationsPaginated(creatorID, search, status, level, idCategory string, limit, offset int) ([]models.Formation, int, error) {
	where := `WHERE f.id_creator = ? AND f.deleted_at IS NULL`
	args := []any{creatorID}

	if search != "" {
		where += " AND f.title_formation LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if status != "" {
		where += " AND f.status = ?"
		args = append(args, status)
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
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM formation f "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT f.id_formation, f.title_formation, f.description_formation, f.date_formation,
		       f.address_formation, f.city_formation, f.postal_formation, f.capacity, f.level, f.duration_hours, f.status,
		       f.rejection_reason, f.id_category,
		       c.name_category,
		       f.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       f.id_formateur, CONCAT(fmt.first_name, ' ', fmt.last_name) AS formateur_name,
		       f.syllabus, f.prerequisites, f.objectives,
		       COUNT(ins.id_user) AS inscription_count,
		       f.created_at, f.price
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
			&f.AddressFormation, &f.CityFormation, &f.PostalFormation, &f.Capacity, &f.Level, &f.DurationHours, &f.Status,
			&f.RejectionReason, &f.IdCategory, &f.CategoryName,
			&f.IdCreator, &f.CreatorName,
			&f.IdFormateur, &f.FormateurName,
			&f.Syllabus, &f.Prerequisites, &f.Objectives,
			&f.InscriptionCount,
			&f.CreatedAt, &f.Price,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, f)
	}
	return list, total, nil
}
