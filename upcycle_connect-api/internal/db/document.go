package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func CreateDocument(doc models.Document) error {
	_, err := config.Conn.Exec(
		"INSERT INTO DOCUMENT (id_document, id_user, category, link) VALUES (?, ?, ?, ?)",
		doc.IdDocument, doc.IdUser, doc.Category, doc.Link,
	)
	return err
}

func GetDocumentById(id string) (models.Document, error) {
	var d models.Document
	err := config.Conn.QueryRow(
		"SELECT id_document, id_user, category, link FROM DOCUMENT WHERE id_document = ?",
		id,
	).Scan(&d.IdDocument, &d.IdUser, &d.Category, &d.Link)
	return d, err
}

func DeleteDocument(id string) error {
	_, err := config.Conn.Exec("DELETE FROM DOCUMENT WHERE id_document = ?", id)
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
		if err := rows.Scan(&d.IdDocument, &d.IdUser, &d.Category, &d.Link); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}
