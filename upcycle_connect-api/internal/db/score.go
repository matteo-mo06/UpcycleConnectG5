package db

import (
	"database/sql"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func AwardScore(userID, actionType, idReference string) error {
	tx, err := config.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var points int
	err = tx.QueryRow(
		"SELECT points FROM score_action WHERE action_type = ?", actionType,
	).Scan(&points)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		"INSERT IGNORE INTO score_log (id, id_user, action_type, id_reference, points) VALUES (?, ?, ?, ?, ?)",
		uuid.New().String(), userID, actionType, idReference, points,
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return tx.Commit()
	}

	_, err = tx.Exec(
		"UPDATE user SET upcycling_score = upcycling_score + ? WHERE id_user = ?",
		points, userID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetScoreActions() ([]models.ScoreAction, error) {
	rows, err := config.Conn.Query(
		"SELECT action_type, points, description FROM score_action ORDER BY points DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.ScoreAction{}
	for rows.Next() {
		var a models.ScoreAction
		if err := rows.Scan(&a.ActionType, &a.Points, &a.Description); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func UpdateScoreAction(actionType string, points int) (bool, error) {
	res, err := config.Conn.Exec(
		"UPDATE score_action SET points = ? WHERE action_type = ?",
		points, actionType,
	)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

func GetUserScoreBreakdown(userID string) (models.ScoreBreakdown, error) {
	var result models.ScoreBreakdown

	err := config.Conn.QueryRow(
		"SELECT upcycling_score FROM USER WHERE id_user = ?", userID,
	).Scan(&result.Score)
	if err != nil {
		return result, err
	}

	rows, err := config.Conn.Query(`
		SELECT sa.action_type, sa.description, sa.points,
		       COUNT(sl.id) AS earned_count,
		       COALESCE(SUM(sl.points), 0) AS subtotal
		FROM score_action sa
		LEFT JOIN score_log sl ON sl.action_type = sa.action_type AND sl.id_user = ?
		GROUP BY sa.action_type, sa.description, sa.points
		ORDER BY sa.points DESC`, userID)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	breakdown := []models.ScoreBreakdownItem{}
	for rows.Next() {
		var item models.ScoreBreakdownItem
		if err := rows.Scan(&item.ActionType, &item.Description, &item.Points, &item.Count, &item.Subtotal); err != nil {
			return result, err
		}
		breakdown = append(breakdown, item)
	}

	result.Breakdown = breakdown
	return result, nil
}
