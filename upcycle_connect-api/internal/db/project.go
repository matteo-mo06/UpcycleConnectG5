package db

import (
	"database/sql"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

var publicStatuses = "p.status IN ('open', 'in_progress', 'completed')"

func IsPublicProject(status string) bool {
	return status == "open" || status == "in_progress" || status == "completed"
}

func GetProjects(userID, search string, limit, offset int) ([]models.Project, int, error) {
	where := "(" + publicStatuses + " OR p.id_creator = ?)"
	countArgs := []any{userID}
	if search != "" {
		where += " AND p.title_project LIKE ?"
		countArgs = append(countArgs, "%"+search+"%")
	}

	var total int
	if err := config.Conn.QueryRow(`SELECT COUNT(*) FROM PROJECT p WHERE `+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := []any{userID, userID, userID}
	if search != "" {
		queryArgs = append(queryArgs, "%"+search+"%")
	}
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT p.id_project, p.title_project, p.description_project, p.start_date_project,
		       p.end_date, p.location_project, p.capacity, p.status, p.rejection_reason,
		       p.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       COUNT(ins.id_user) AS members_count,
		       EXISTS(SELECT 1 FROM USER_PROJECT_INSCRIPTION r WHERE r.id_user = ? AND r.id_project = p.id_project) AS is_registered,
		       p.created_at,
		       (SELECT COUNT(*) FROM USER_PROJECT_UPDOWN WHERE id_project = p.id_project AND updown = 1) AS up_votes,
		       (SELECT COUNT(*) FROM USER_PROJECT_UPDOWN WHERE id_project = p.id_project AND updown = -1) AS down_votes,
		       COALESCE((SELECT updown FROM USER_PROJECT_UPDOWN WHERE id_project = p.id_project AND id_user = ?), 0) AS user_vote
		FROM PROJECT p
		LEFT JOIN USER u ON u.id_user = p.id_creator
		LEFT JOIN USER_PROJECT_INSCRIPTION ins ON ins.id_project = p.id_project
		WHERE `+where+`
		GROUP BY p.id_project
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Project
	for rows.Next() {
		var p models.Project
		var rejReason sql.NullString
		err := rows.Scan(
			&p.IdProject, &p.TitleProject, &p.DescriptionProject, &p.StartDateProject,
			&p.EndDate, &p.LocationProject, &p.Capacity, &p.Status, &rejReason,
			&p.IdCreator, &p.CreatorName, &p.MembersCount, &p.IsRegistered,
			&p.CreatedAt, &p.UpVotes, &p.DownVotes, &p.UserVote,
		)
		if err != nil {
			return nil, 0, err
		}
		if rejReason.Valid {
			p.RejectionReason = &rejReason.String
		}
		list = append(list, p)
	}
	return list, total, nil
}

func GetAllProjectsAdmin(search, status string, limit, offset int) ([]models.Project, int, error) {
	where := "WHERE 1=1"
	var args []any

	if search != "" {
		where += " AND p.title_project LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if status != "" {
		where += " AND p.status = ?"
		args = append(args, status)
	}

	countArgs := make([]any, len(args))
	copy(countArgs, args)

	var total int
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM PROJECT p "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT p.id_project, p.title_project, p.description_project, p.start_date_project,
		       p.end_date, p.location_project, p.capacity, p.status, p.rejection_reason,
		       p.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       COUNT(ins.id_user) AS members_count,
		       p.created_at,
		       (SELECT COUNT(*) FROM USER_PROJECT_UPDOWN WHERE id_project = p.id_project AND updown = 1) AS up_votes,
		       (SELECT COUNT(*) FROM USER_PROJECT_UPDOWN WHERE id_project = p.id_project AND updown = -1) AS down_votes
		FROM PROJECT p
		LEFT JOIN USER u ON u.id_user = p.id_creator
		LEFT JOIN USER_PROJECT_INSCRIPTION ins ON ins.id_project = p.id_project
		`+where+`
		GROUP BY p.id_project
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Project
	for rows.Next() {
		var p models.Project
		var rejReason sql.NullString
		err := rows.Scan(
			&p.IdProject, &p.TitleProject, &p.DescriptionProject, &p.StartDateProject,
			&p.EndDate, &p.LocationProject, &p.Capacity, &p.Status, &rejReason,
			&p.IdCreator, &p.CreatorName, &p.MembersCount,
			&p.CreatedAt, &p.UpVotes, &p.DownVotes,
		)
		if err != nil {
			return nil, 0, err
		}
		if rejReason.Valid {
			p.RejectionReason = &rejReason.String
		}
		list = append(list, p)
	}
	return list, total, nil
}

func GetProjectById(id string) (models.Project, error) {
	var p models.Project
	var rejReason sql.NullString
	err := config.Conn.QueryRow(`
		SELECT p.id_project, p.title_project, p.description_project, p.start_date_project,
		       p.end_date, p.location_project, p.capacity, p.status, p.rejection_reason,
		       p.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       (COUNT(ins.id_user) + 1) AS members_count,
		       p.created_at
		FROM PROJECT p
		LEFT JOIN USER u ON u.id_user = p.id_creator
		LEFT JOIN USER_PROJECT_INSCRIPTION ins ON ins.id_project = p.id_project
		WHERE p.id_project = ?
		GROUP BY p.id_project`, id,
	).Scan(
		&p.IdProject, &p.TitleProject, &p.DescriptionProject, &p.StartDateProject,
		&p.EndDate, &p.LocationProject, &p.Capacity, &p.Status, &rejReason,
		&p.IdCreator, &p.CreatorName, &p.MembersCount,
		&p.CreatedAt,
	)
	if rejReason.Valid {
		p.RejectionReason = &rejReason.String
	}
	return p, err
}

func GetProjectOwnerID(projectID string) (string, error) {
	var ownerID sql.NullString
	err := config.Conn.QueryRow(
		`SELECT id_creator FROM PROJECT WHERE id_project = ?`, projectID,
	).Scan(&ownerID)
	if err != nil {
		return "", err
	}
	return ownerID.String, nil
}

func CreateProject(p models.Project) error {
	_, err := config.Conn.Exec(`
		INSERT INTO PROJECT (
			id_project, title_project, description_project, start_date_project,
			end_date, location_project, capacity, status, id_creator, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', ?, NOW())`,
		p.IdProject, p.TitleProject, p.DescriptionProject, p.StartDateProject,
		p.EndDate, p.LocationProject, p.Capacity, p.IdCreator,
	)
	return err
}

func UpdateProject(p models.Project) error {
	_, err := config.Conn.Exec(`
		UPDATE PROJECT SET
			title_project = ?,
			description_project = ?,
			start_date_project = ?,
			end_date = ?,
			location_project = ?,
			capacity = ?,
			status = 'pending',
			rejection_reason = NULL
		WHERE id_project = ?`,
		p.TitleProject, p.DescriptionProject, p.StartDateProject,
		p.EndDate, p.LocationProject, p.Capacity, p.IdProject,
	)
	return err
}

func UpdateProjectAdmin(p models.Project) error {
	_, err := config.Conn.Exec(`
		UPDATE PROJECT SET
			title_project = ?,
			description_project = ?,
			start_date_project = ?,
			end_date = ?,
			location_project = ?,
			capacity = ?,
			status = ?
		WHERE id_project = ?`,
		p.TitleProject, p.DescriptionProject, p.StartDateProject,
		p.EndDate, p.LocationProject, p.Capacity, p.Status, p.IdProject,
	)
	return err
}

func ApproveProject(id string) error {
	_, err := config.Conn.Exec(
		"UPDATE PROJECT SET status = 'open', rejection_reason = NULL WHERE id_project = ?", id,
	)
	return err
}

func RejectProject(id string, reason *string) error {
	_, err := config.Conn.Exec(
		"UPDATE PROJECT SET status = 'rejected', rejection_reason = ? WHERE id_project = ?", reason, id,
	)
	return err
}

func DeleteProject(id string) error {
	if _, err := config.Conn.Exec("DELETE FROM USER_PROJECT_INSCRIPTION WHERE id_project = ?", id); err != nil {
		return err
	}
	if _, err := config.Conn.Exec("DELETE FROM USER_PROJECT_UPDOWN WHERE id_project = ?", id); err != nil {
		return err
	}
	if _, err := config.Conn.Exec("DELETE FROM PAYEMENT_PROJECT WHERE id_project = ?", id); err != nil {
		return err
	}
	_, err := config.Conn.Exec("DELETE FROM PROJECT WHERE id_project = ?", id)
	return err
}

func RegisterUserForProject(userID, projectID string) error {
	_, err := config.Conn.Exec(
		"INSERT IGNORE INTO USER_PROJECT_INSCRIPTION (id_user, id_project) VALUES (?, ?)",
		userID, projectID,
	)
	return err
}

func UnregisterUserFromProject(userID, projectID string) error {
	_, err := config.Conn.Exec(
		"DELETE FROM USER_PROJECT_INSCRIPTION WHERE id_user = ? AND id_project = ?",
		userID, projectID,
	)
	return err
}

func VoteProject(userID, projectID string, value int) error {
	_, err := config.Conn.Exec(
		`INSERT INTO USER_PROJECT_UPDOWN (id_user, id_project, updown) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE updown = ?`,
		userID, projectID, value, value,
	)
	return err
}

func RemoveProjectVote(userID, projectID string) error {
	_, err := config.Conn.Exec(
		`DELETE FROM USER_PROJECT_UPDOWN WHERE id_user = ? AND id_project = ?`,
		userID, projectID,
	)
	return err
}

func GetMyCreatedProjects(creatorID string) ([]models.Project, error) {
	rows, err := config.Conn.Query(`
		SELECT p.id_project, p.title_project, p.description_project, p.start_date_project,
		       p.end_date, p.location_project, p.capacity, p.status, p.rejection_reason,
		       p.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       COUNT(ins.id_user) AS members_count,
		       p.created_at,
		       (SELECT COUNT(*) FROM USER_PROJECT_UPDOWN WHERE id_project = p.id_project AND updown = 1) AS up_votes,
		       (SELECT COUNT(*) FROM USER_PROJECT_UPDOWN WHERE id_project = p.id_project AND updown = -1) AS down_votes
		FROM PROJECT p
		LEFT JOIN USER u ON u.id_user = p.id_creator
		LEFT JOIN USER_PROJECT_INSCRIPTION ins ON ins.id_project = p.id_project
		WHERE p.id_creator = ?
		GROUP BY p.id_project
		ORDER BY p.created_at DESC`, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Project
	for rows.Next() {
		var p models.Project
		var rejReason sql.NullString
		err := rows.Scan(
			&p.IdProject, &p.TitleProject, &p.DescriptionProject, &p.StartDateProject,
			&p.EndDate, &p.LocationProject, &p.Capacity, &p.Status, &rejReason,
			&p.IdCreator, &p.CreatorName, &p.MembersCount,
			&p.CreatedAt, &p.UpVotes, &p.DownVotes,
		)
		if err != nil {
			return nil, err
		}
		if rejReason.Valid {
			p.RejectionReason = &rejReason.String
		}
		list = append(list, p)
	}
	return list, nil
}

func GetProjectStatusCounts() (models.ProjectStats, error) {
	var s models.ProjectStats
	err := config.Conn.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(status = 'pending'), 0),
			COALESCE(SUM(status = 'open'), 0),
			COALESCE(SUM(status = 'in_progress'), 0),
			COALESCE(SUM(status = 'completed'), 0),
			COALESCE(SUM(status = 'rejected'), 0)
		FROM PROJECT
	`).Scan(&s.Total, &s.Pending, &s.Open, &s.InProgress, &s.Completed, &s.Rejected)
	return s, err
}

func GetUserRegisteredProjects(userID string) ([]models.Project, error) {
	rows, err := config.Conn.Query(`
		SELECT p.id_project, p.title_project, p.description_project, p.start_date_project,
		       p.end_date, p.location_project, p.capacity, p.status,
		       p.id_creator, CONCAT(u.first_name, ' ', u.last_name) AS creator_name,
		       COUNT(ins.id_user) AS members_count,
		       p.created_at
		FROM PROJECT p
		LEFT JOIN USER u ON u.id_user = p.id_creator
		LEFT JOIN USER_PROJECT_INSCRIPTION ins ON ins.id_project = p.id_project
		WHERE EXISTS (
			SELECT 1 FROM USER_PROJECT_INSCRIPTION WHERE id_user = ? AND id_project = p.id_project
		)
		GROUP BY p.id_project
		ORDER BY p.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Project
	for rows.Next() {
		var p models.Project
		p.IsRegistered = true
		err := rows.Scan(
			&p.IdProject, &p.TitleProject, &p.DescriptionProject, &p.StartDateProject,
			&p.EndDate, &p.LocationProject, &p.Capacity, &p.Status,
			&p.IdCreator, &p.CreatorName, &p.MembersCount,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
