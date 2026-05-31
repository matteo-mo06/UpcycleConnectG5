package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func CreateDocument(doc models.Document) error {
	_, err := config.Conn.Exec(
		"INSERT INTO DOCUMENT (id_document, id_user, category, link) VALUES (?, ?, ?, ?)",
		doc.Id_document, doc.Id_user, doc.Category, doc.Link,
	)
	return err
}

func GetDocumentsByCategory(category string) ([]models.Document, error) {
	rows, err := config.Conn.Query(
		"SELECT id_document, id_user, category, link FROM DOCUMENT WHERE category = ?",
		category,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.Document
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(&d.Id_document, &d.Id_user, &d.Category, &d.Link); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}
