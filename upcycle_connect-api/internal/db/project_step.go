package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetProjectSteps(projectID string) ([]models.ProjectStep, error) {
	rows, err := config.Conn.Query(
		`SELECT id_step, id_project, title, description, status, step_order FROM PROJECT_STEP WHERE id_project = ? ORDER BY step_order, id_step`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ProjectStep
	for rows.Next() {
		var s models.ProjectStep
		if err := rows.Scan(&s.IdStep, &s.IdProject, &s.Title, &s.Description, &s.Status, &s.StepOrder); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func CreateProjectStep(s models.ProjectStep) error {
	_, err := config.Conn.Exec(
		`INSERT INTO PROJECT_STEP (id_step, id_project, title, description, status, step_order) VALUES (?, ?, ?, ?, ?, ?)`,
		s.IdStep, s.IdProject, s.Title, s.Description, s.Status, s.StepOrder,
	)
	return err
}

func UpdateProjectStep(s models.ProjectStep) error {
	_, err := config.Conn.Exec(
		`UPDATE PROJECT_STEP SET title = ?, description = ?, step_order = ? WHERE id_step = ?`,
		s.Title, s.Description, s.StepOrder, s.IdStep,
	)
	return err
}

func UpdateProjectStepStatus(id, status string) error {
	_, err := config.Conn.Exec(
		`UPDATE PROJECT_STEP SET status = ? WHERE id_step = ?`, status, id,
	)
	return err
}

func DeleteProjectStep(id string) error {
	_, err := config.Conn.Exec(`DELETE FROM PROJECT_STEP WHERE id_step = ?`, id)
	return err
}
