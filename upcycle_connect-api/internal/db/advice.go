package db

import (
	"database/sql"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAdvices(search, idCategory string, limit, offset int) ([]models.Advice, int, error) {
	where := "WHERE 1=1"
	var args []any

	if search != "" {
		where += " AND a.title LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if idCategory != "" {
		where += " AND a.id_category = ?"
		args = append(args, idCategory)
	}

	countArgs := make([]any, len(args))
	copy(countArgs, args)
	var total int
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM advice a "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT a.id_advice, a.title, a.id_category, c.name_category,
		       a.description, a.date_advice,
		       a.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       a.created_at
		FROM advice a
		LEFT JOIN category c ON c.id_category = a.id_category
		LEFT JOIN user u ON u.id_user = a.id_creator
		`+where+`
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Advice
	for rows.Next() {
		var a models.Advice
		if err := rows.Scan(
			&a.IdAdvice, &a.Title, &a.IdCategory, &a.CategoryName,
			&a.Description, &a.DateAdvice,
			&a.IdCreator, &a.CreatorName,
			&a.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}
	return list, total, nil
}

func GetAdviceById(id string) (models.Advice, error) {
	var a models.Advice
	err := config.Conn.QueryRow(`
		SELECT a.id_advice, a.title, a.id_category, c.name_category,
		       a.description, a.date_advice,
		       a.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       a.created_at
		FROM advice a
		LEFT JOIN category c ON c.id_category = a.id_category
		LEFT JOIN user u ON u.id_user = a.id_creator
		WHERE a.id_advice = ?`, id,
	).Scan(
		&a.IdAdvice, &a.Title, &a.IdCategory, &a.CategoryName,
		&a.Description, &a.DateAdvice,
		&a.IdCreator, &a.CreatorName,
		&a.CreatedAt,
	)
	return a, err
}

func CreateAdvice(creatorID, title string, description, idCategory, dateAdvice *string, photoURLs []string) error {
	id := uuid.New().String()
	_, err := config.Conn.Exec(`
		INSERT INTO advice (id_advice, title, description, id_category, date_advice, id_creator, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())`,
		id, title, description, idCategory, dateAdvice, creatorID,
	)
	if err != nil {
		return err
	}
	for _, url := range photoURLs {
		if _, err := config.Conn.Exec(
			"INSERT INTO DOCUMENT (id_document, id_user, category, link) VALUES (?, ?, ?, ?)",
			uuid.New().String(), creatorID, id, url,
		); err != nil {
			return err
		}
	}
	return nil
}

func UpdateAdvice(id, title string, description, idCategory, dateAdvice *string) error {
	_, err := config.Conn.Exec(`
		UPDATE advice SET title = ?, description = ?, id_category = ?, date_advice = ?
		WHERE id_advice = ?`,
		title, description, idCategory, dateAdvice, id,
	)
	return err
}

func DeleteAdvice(id string) error {
	_, err := config.Conn.Exec("DELETE FROM advice WHERE id_advice = ?", id)
	return err
}

func GetAdviceCreatorID(id string) (string, error) {
	var creatorID sql.NullString
	err := config.Conn.QueryRow(
		"SELECT id_creator FROM advice WHERE id_advice = ?", id,
	).Scan(&creatorID)
	if err != nil {
		return "", err
	}
	return creatorID.String, nil
}
