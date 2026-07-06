package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetFormationSteps(formationID string) ([]models.FormationStep, error) {
	rows, err := config.Conn.Query(
		`SELECT id_step, id_formation, title, description, status, step_order FROM formation_step WHERE id_formation = ? ORDER BY step_order, id_step`,
		formationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.FormationStep
	for rows.Next() {
		var s models.FormationStep
		if err := rows.Scan(&s.IdStep, &s.IdFormation, &s.Title, &s.Description, &s.Status, &s.StepOrder); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func CreateFormationStep(s models.FormationStep) error {
	_, err := config.Conn.Exec(
		`INSERT INTO formation_step (id_step, id_formation, title, description, status, step_order) VALUES (?, ?, ?, ?, ?, ?)`,
		s.IdStep, s.IdFormation, s.Title, s.Description, s.Status, s.StepOrder,
	)
	return err
}

func UpdateFormationStep(s models.FormationStep) error {
	_, err := config.Conn.Exec(
		`UPDATE formation_step SET title = ?, description = ?, step_order = ? WHERE id_step = ?`,
		s.Title, s.Description, s.StepOrder, s.IdStep,
	)
	return err
}

func UpdateFormationStepStatus(id, status string) error {
	_, err := config.Conn.Exec(
		`UPDATE formation_step SET status = ? WHERE id_step = ?`, status, id,
	)
	return err
}

func DeleteFormationStep(id string) error {
	_, err := config.Conn.Exec(`DELETE FROM formation_step WHERE id_step = ?`, id)
	return err
}
