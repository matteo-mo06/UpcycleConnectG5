package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllCategories() ([]models.Category, error) {
	rows, err := config.Conn.Query(`SELECT id_category, name_category, description_category FROM CATEGORY WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.IdCategory, &c.NameCategory, &c.DescriptionCategory)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func GetCategoryById(id string) (models.Category, error) {
	var c models.Category
	err := config.Conn.QueryRow(
		`SELECT id_category, name_category, description_category FROM CATEGORY WHERE id_category = ? AND deleted_at IS NULL`, id,
	).Scan(&c.IdCategory, &c.NameCategory, &c.DescriptionCategory)
	return c, err
}

func CreateCategory(c models.Category) error {
	_, err := config.Conn.Exec(`
		INSERT INTO CATEGORY (
			id_category,
			name_category,
			description_category
		)
		VALUES (?, ?, ?)
	`,
		c.IdCategory,
		c.NameCategory,
		c.DescriptionCategory,
	)
	return err
}

func UpdateCategory(c models.Category) error {
	_, err := config.Conn.Exec(`
		UPDATE CATEGORY SET
			name_category = ?,
			description_category = ?
		WHERE id_category = ?
	`,
		c.NameCategory,
		c.DescriptionCategory,
		c.IdCategory,
	)
	return err
}

func DeleteCategory(id string) error {
	_, err := config.Conn.Exec("UPDATE CATEGORY SET deleted_at = NOW() WHERE id_category = ?", id)
	return err
}
