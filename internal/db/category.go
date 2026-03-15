package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllCategories() ([]models.Category, error) {
	rows, err := config.Conn.Query(`SELECT id_category, name_category, description_category FROM CATEGORY`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.Id_category, &c.Name_category, &c.Description_category); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func GetCategoryById(id string) (models.Category, error) {
	var c models.Category
	err := config.Conn.QueryRow(
		`SELECT id_category, name_category, description_category FROM CATEGORY WHERE id_category = ?`, id,
	).Scan(&c.Id_category, &c.Name_category, &c.Description_category)
	return c, err
}

func CreateCategory(c models.Category) error {
	_, err := config.Conn.Exec(
		`INSERT INTO CATEGORY (id_category, name_category, description_category) VALUES (?, ?, ?)`,
		c.Id_category, c.Name_category, c.Description_category,
	)
	return err
}

func UpdateCategory(c models.Category) error {
	_, err := config.Conn.Exec(
		`UPDATE CATEGORY SET name_category = ?, description_category = ? WHERE id_category = ?`,
		c.Name_category, c.Description_category, c.Id_category,
	)
	return err
}

func DeleteCategory(id string) error {
	_, err := config.Conn.Exec("UPDATE ANNOUNCEMENT SET id_category = NULL WHERE id_category = ?", id)
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec("DELETE FROM CATEGORY WHERE id_category = ?", id)
	return err
}
