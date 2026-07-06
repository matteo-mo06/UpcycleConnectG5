package db

import (
	"database/sql"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetProfessionalRequests(status string) ([]models.ProfessionalRequestDetail, error) {
	query := `
		SELECT pr.id_request, pr.id_user, pr.status, pr.created_at, pr.processed_at, u.email, u.first_name, u.last_name
		FROM PROFESSIONAL_REQUEST pr
		JOIN USER u ON u.id_user = pr.id_user
		WHERE 1=1`

	var args []any
	if status != "" {
		query += " AND pr.status = ?"
		args = append(args, status)
	}
	query += " ORDER BY pr.created_at DESC"

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ProfessionalRequestDetail
	for rows.Next() {
		var req models.ProfessionalRequestDetail
		err := rows.Scan(
			&req.IdRequest,
			&req.IdUser,
			&req.Status,
			&req.CreatedAt,
			&req.ProcessedAt,
			&req.Email,
			&req.FirstName,
			&req.LastName,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func GetProfessionalRequestById(id string) (models.ProfessionalRequest, error) {
	var req models.ProfessionalRequest
	row := config.Conn.QueryRow(`
		SELECT id_request, id_user, status, created_at, processed_at
		FROM PROFESSIONAL_REQUEST
		WHERE id_request = ?`, id)

	err := row.Scan(&req.IdRequest, &req.IdUser, &req.Status, &req.CreatedAt, &req.ProcessedAt)
	return req, err
}

func CreateProfessionalRequest(req models.ProfessionalRequest) error {
	_, err := config.Conn.Exec(`
		INSERT INTO PROFESSIONAL_REQUEST (id_request, id_user, status)
		VALUES (?, ?, 'pending')`,
		req.IdRequest, req.IdUser,
	)
	return err
}

func UpdateProfessionalRequestStatus(id, status string) error {
	_, err := config.Conn.Exec(`
		UPDATE PROFESSIONAL_REQUEST
		SET status = ?, processed_at = NOW()
		WHERE id_request = ?`,
		status, id,
	)
	return err
}

func GetPendingRequestByUser(userID string) (models.ProfessionalRequest, error) {
	var req models.ProfessionalRequest
	row := config.Conn.QueryRow(`
		SELECT id_request, id_user, status, created_at, processed_at
		FROM PROFESSIONAL_REQUEST
		WHERE id_user = ? AND status = 'pending'`, userID)

	err := row.Scan(&req.IdRequest, &req.IdUser, &req.Status, &req.CreatedAt, &req.ProcessedAt)
	if err == sql.ErrNoRows {
		return req, sql.ErrNoRows
	}
	return req, err
}

func GetRequestsByUser(userID string) ([]models.ProfessionalRequest, error) {
	rows, err := config.Conn.Query(`
		SELECT id_request, id_user, status, created_at, processed_at
		FROM PROFESSIONAL_REQUEST
		WHERE id_user = ?
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ProfessionalRequest
	for rows.Next() {
		var req models.ProfessionalRequest
		if err := rows.Scan(&req.IdRequest, &req.IdUser, &req.Status, &req.CreatedAt, &req.ProcessedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func GetLatestRequestByUser(userID string) (models.ProfessionalRequest, error) {
	var req models.ProfessionalRequest
	row := config.Conn.QueryRow(`
		SELECT id_request, id_user, status, created_at, processed_at
		FROM PROFESSIONAL_REQUEST
		WHERE id_user = ?
		ORDER BY created_at DESC
		LIMIT 1`, userID)

	err := row.Scan(&req.IdRequest, &req.IdUser, &req.Status, &req.CreatedAt, &req.ProcessedAt)
	if err == sql.ErrNoRows {
		return req, sql.ErrNoRows
	}
	return req, err
}
