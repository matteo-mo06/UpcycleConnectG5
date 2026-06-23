package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetProjectMaterials(projectID string) ([]models.ProjectMaterial, error) {
	rows, err := config.Conn.Query(
		`SELECT id_material, id_project, name, quantity, unit FROM PROJECT_MATERIAL WHERE id_project = ? ORDER BY name`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ProjectMaterial
	for rows.Next() {
		var m models.ProjectMaterial
		if err := rows.Scan(&m.IdMaterial, &m.IdProject, &m.Name, &m.Quantity, &m.Unit); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func CreateProjectMaterial(m models.ProjectMaterial) error {
	_, err := config.Conn.Exec(
		`INSERT INTO PROJECT_MATERIAL (id_material, id_project, name, quantity, unit) VALUES (?, ?, ?, ?, ?)`,
		m.IdMaterial, m.IdProject, m.Name, m.Quantity, m.Unit,
	)
	return err
}

func UpdateProjectMaterial(m models.ProjectMaterial) error {
	_, err := config.Conn.Exec(
		`UPDATE PROJECT_MATERIAL SET name = ?, quantity = ?, unit = ? WHERE id_material = ?`,
		m.Name, m.Quantity, m.Unit, m.IdMaterial,
	)
	return err
}

func DeleteProjectMaterial(id string) error {
	_, err := config.Conn.Exec(`DELETE FROM PROJECT_MATERIAL WHERE id_material = ?`, id)
	return err
}
