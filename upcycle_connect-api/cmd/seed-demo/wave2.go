package main

import (
	"database/sql"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

var stepTitles = []string{"Préparation du matériel", "Découpe et assemblage", "Finitions", "Présentation finale"}
var stepStatuses = []string{"todo", "in_progress", "done"}
var materialNames = []string{"Bois", "Vis", "Peinture", "Tissu", "Colle", "Clous", "Vernis", "Ficelle"}
var materialUnits = []string{"kg", "m", "pièce(s)", "L"}

func nonDeletedAnnouncements(s *state) []string {
	var out []string
	for _, id := range s.announcementIDs {
		if !s.announcementDeleted[id] {
			out = append(out, id)
		}
	}
	return out
}

func seedAnnouncementLocker(db *sql.DB, s *state, r *refs) error {
	eligible := nonDeletedAnnouncements(s)
	count := 0
	for i, lockerID := range r.lockerIDs {
		if i >= len(eligible) {
			break
		}
		annID := eligible[i]
		if _, err := db.Exec(
			"INSERT INTO announcement_locker (id_announcement, id_locker) VALUES (?, ?)", annID, lockerID,
		); err != nil {
			return err
		}
		s.announcementLockerIDs = append(s.announcementLockerIDs, annID)
		count++
	}
	fmt.Printf("announcement_locker: %d créés\n", count)
	return nil
}

func seedPairTable(db *sql.DB, table, colA, colB string, pairs [][2]string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)", table, colA, colB)
	for _, p := range pairs {
		if _, err := db.Exec(query, p[0], p[1]); err != nil {
			return err
		}
	}
	return nil
}

func seedUserProjectInscription(db *sql.DB, s *state) ([][2]string, error) {
	pairs := fillPairs(90, s.userIDs, s.projectIDs, nil)
	if err := seedPairTable(db, "user_project_inscription", "id_user", "id_project", pairs); err != nil {
		return nil, err
	}
	fmt.Printf("user_project_inscription: %d créés\n", len(pairs))
	return pairs, nil
}

func seedUserProjectUpdown(db *sql.DB, s *state) error {
	pairs := fillPairs(150, s.userIDs, s.projectIDs, nil)
	for _, p := range pairs {
		updown := 1
		if percent(30) {
			updown = -1
		}
		if _, err := db.Exec(
			"INSERT INTO user_project_updown (id_user, id_project, updown) VALUES (?, ?, ?)", p[0], p[1], updown,
		); err != nil {
			return err
		}
	}
	fmt.Printf("user_project_updown: %d créés\n", len(pairs))
	return nil
}

func seedUserFormationInscription(db *sql.DB, s *state) ([][2]string, error) {
	pairs := fillPairs(110, s.userIDs, s.formationIDs, nil)
	if err := seedPairTable(db, "user_formation_inscription", "id_user", "id_formation", pairs); err != nil {
		return nil, err
	}
	fmt.Printf("user_formation_inscription: %d créés\n", len(pairs))
	return pairs, nil
}

func seedUserEventInscription(db *sql.DB, s *state) ([][2]string, error) {
	pairs := fillPairs(110, s.userIDs, s.eventIDs, nil)
	if err := seedPairTable(db, "user_event_inscription", "id_user", "id_event", pairs); err != nil {
		return nil, err
	}
	fmt.Printf("user_event_inscription: %d créés\n", len(pairs))
	return pairs, nil
}

func seedUserAnnouncement(db *sql.DB, s *state) error {
	for _, annID := range s.announcementIDs {
		owner := s.announcementOwner[annID]
		if _, err := db.Exec(
			"INSERT INTO user_announcement (id_user, id_announcement) VALUES (?, ?)", owner, annID,
		); err != nil {
			return err
		}
	}
	fmt.Printf("user_announcement: %d créés (1 par annonce, propriétaire)\n", len(s.announcementIDs))
	return nil
}

func seedUserTopic(db *sql.DB, s *state) error {
	seen := newPairSet()
	count := 0
	for _, topicID := range s.topicIDs {
		author := s.topicAuthor[topicID]
		if seen.add(author, topicID) {
			if _, err := db.Exec(
				"INSERT INTO user_topic (id_user, id_topic) VALUES (?, ?)", author, topicID,
			); err != nil {
				return err
			}
			count++
		}
	}
	extra := fillPairs(110-count, s.userIDs, s.topicIDs, seen)
	if err := seedPairTable(db, "user_topic", "id_user", "id_topic", extra); err != nil {
		return err
	}
	fmt.Printf("user_topic: %d créés\n", count+len(extra))
	return nil
}

func seedUserAdvice(db *sql.DB, s *state) error {
	pairs := fillPairs(110, s.userIDs, s.adviceIDs, nil)
	if err := seedPairTable(db, "user_advice", "id_user", "id_advice", pairs); err != nil {
		return err
	}
	fmt.Printf("user_advice: %d créés\n", len(pairs))
	return nil
}

func seedTopicPost(db *sql.DB, s *state) error {
	for _, postID := range s.postIDs {
		topicID := s.postTopic[postID]
		if _, err := db.Exec(
			"INSERT INTO topic_post (id_topic, id_post) VALUES (?, ?)", topicID, postID,
		); err != nil {
			return err
		}
	}
	fmt.Printf("topic_post: %d créés\n", len(s.postIDs))
	return nil
}

func seedProjectMaterialAndSteps(db *sql.DB, s *state) error {
	materialCount, stepCount := 0, 0
	for _, projectID := range s.projectIDs {
		for i := 0; i < 3; i++ {
			if _, err := db.Exec(
				"INSERT INTO project_material (id_material, id_project, name, quantity, unit) VALUES (?, ?, ?, ?, ?)",
				newID(), projectID, pick(materialNames), randMoney(1, 20), pick(materialUnits),
			); err != nil {
				return err
			}
			materialCount++

			if _, err := db.Exec(
				"INSERT INTO project_step (id_step, id_project, title, description, status, step_order) VALUES (?, ?, ?, ?, ?, ?)",
				newID(), projectID, stepTitles[i%len(stepTitles)], "Étape : "+stepTitles[i%len(stepTitles)], pick(stepStatuses), i+1,
			); err != nil {
				return err
			}
			stepCount++
		}
	}
	fmt.Printf("project_material: %d créés, project_step: %d créés\n", materialCount, stepCount)
	return nil
}

func seedFormationSteps(db *sql.DB, s *state) error {
	count := 0
	for _, formationID := range s.formationIDs {
		for i := 0; i < 3; i++ {
			if _, err := db.Exec(
				"INSERT INTO formation_step (id_step, id_formation, title, description, status, step_order) VALUES (?, ?, ?, ?, ?, ?)",
				newID(), formationID, stepTitles[i%len(stepTitles)], "Étape : "+stepTitles[i%len(stepTitles)], pick(stepStatuses), i+1,
			); err != nil {
				return err
			}
			count++
		}
	}
	fmt.Printf("formation_step: %d créés\n", count)
	return nil
}

func seedFeaturedSlotLog(db *sql.DB, s *state) error {
	eligible := nonDeletedAnnouncements(s)
	if len(eligible) > len(s.announcementLockerIDs) {
		eligible = eligible[len(s.announcementLockerIDs):]
	}
	const total = 20
	for i := 0; i < total && i < len(eligible); i++ {
		annID := eligible[i]
		started := daysAgo(randInt(1, 10))
		ends := started.AddDate(0, 0, 7)

		if _, err := db.Exec(
			"UPDATE announcement SET is_featured = 1, featured_until = ?, featured_requested_at = ? WHERE id_announcement = ?",
			sqlDateTime(ends), sqlDateTime(started), annID,
		); err != nil {
			return err
		}
		if _, err := db.Exec(
			"INSERT INTO featured_slot_log (id_slot, id_announcement, started_at, ends_at) VALUES (?, ?, ?, ?)",
			newID(), annID, sqlDateTime(started), sqlDateTime(ends),
		); err != nil {
			return err
		}
	}
	fmt.Printf("featured_slot_log: %d créés\n", total)
	return nil
}

type reportedTarget struct {
	kind         string
	id           string
	reportedUser string
}

func seedReports(db *sql.DB, s *state, r *refs) error {
	var targets []reportedTarget
	for _, id := range s.announcementIDs {
		targets = append(targets, reportedTarget{"announcement", id, s.announcementOwner[id]})
	}
	for _, id := range s.topicIDs {
		targets = append(targets, reportedTarget{"topic", id, s.topicAuthor[id]})
	}
	for _, id := range s.postIDs {
		targets = append(targets, reportedTarget{"post", id, s.postAuthor[id]})
	}
	for _, id := range s.projectIDs {
		targets = append(targets, reportedTarget{"project", id, s.projectCreator[id]})
	}

	const total = 55
	const resolvedCount = 30
	for i := 0; i < total; i++ {
		t := pick(targets)
		id := newID()
		reporter := pickExcept(s.userIDs, t.reportedUser)
		createdAt := gofakeit.DateRange(daysAgo(90), daysAgo(1))

		var idAnnouncement, idTopic, idPost, idProject string
		switch t.kind {
		case "announcement":
			idAnnouncement = t.id
		case "topic":
			idTopic = t.id
		case "post":
			idPost = t.id
		case "project":
			idProject = t.id
		}

		status := "pending"
		resolvedBy := ""
		actionTaken := ""
		resolvedAt := ""
		if i < resolvedCount {
			status = "resolved"
			resolvedBy = r.adminUserID
			actionTaken = pick(reportActionsTaken)
			resolvedAt = sqlDateTime(createdAt.AddDate(0, 0, randInt(1, 5)))
		}

		_, err := db.Exec(`
			INSERT INTO report (id_report, id_user, id_announcement, id_topic, id_post, id_project, reason, status,
				created_at, id_reported_user, resolved_by, action_taken, resolved_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, reporter, nullStr(idAnnouncement), nullStr(idTopic), nullStr(idPost), nullStr(idProject),
			pick(reportReasons), status, sqlDateTime(createdAt), t.reportedUser, nullStr(resolvedBy), nullStr(actionTaken), nullStr(resolvedAt),
		)
		if err != nil {
			return err
		}
		s.reportIDs = append(s.reportIDs, id)
		s.reportedUserFor[id] = t.reportedUser
		if status == "resolved" {
			s.reportResolved = append(s.reportResolved, id)
		}
	}
	fmt.Printf("report: %d créés (%d résolus, %d en attente)\n", total, resolvedCount, total-resolvedCount)
	return nil
}

func loadScoreActionPoints(db *sql.DB) (map[string]int, error) {
	rows, err := db.Query("SELECT action_type, points FROM score_action")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	points := make(map[string]int)
	for rows.Next() {
		var actionType string
		var pts int
		if err := rows.Scan(&actionType, &pts); err != nil {
			return nil, err
		}
		points[actionType] = pts
	}
	return points, nil
}

func seedScoreLog(db *sql.DB, s *state, points map[string]int, projectPairs, formationPairs, eventPairs [][2]string) error {
	seen := make(map[string]bool)
	count := 0
	insert := func(userID, actionType, idReference string) error {
		key := userID + "|" + actionType + "|" + idReference
		if seen[key] {
			return nil
		}
		seen[key] = true
		if _, err := db.Exec(
			"INSERT INTO score_log (id, id_user, action_type, id_reference, points) VALUES (?, ?, ?, ?, ?)",
			newID(), userID, actionType, idReference, points[actionType],
		); err != nil {
			return err
		}
		count++
		return nil
	}

	for _, projectID := range s.projectIDs {
		if err := insert(s.projectCreator[projectID], "project_created", projectID); err != nil {
			return err
		}
	}
	for _, p := range projectPairs {
		if err := insert(p[0], "project_registration", p[1]); err != nil {
			return err
		}
	}
	for _, p := range formationPairs {
		if err := insert(p[0], "formation_registration", p[1]); err != nil {
			return err
		}
	}
	for _, p := range eventPairs {
		if err := insert(p[0], "event_registration", p[1]); err != nil {
			return err
		}
	}
	for _, annID := range s.announcementIDs {
		buyer, ok := s.announcementBuyer[annID]
		if !ok {
			continue
		}
		if err := insert(buyer, "announcement_bought", annID); err != nil {
			return err
		}
		if err := insert(s.announcementOwner[annID], "announcement_sold", annID); err != nil {
			return err
		}
	}
	depositCount := len(s.announcementLockerIDs) * 3 / 4
	for i := 0; i < depositCount; i++ {
		annID := s.announcementLockerIDs[i]
		if err := insert(s.announcementOwner[annID], "deposit_validated", annID); err != nil {
			return err
		}
	}

	fmt.Printf("score_log: %d créés (hors sanctions, ajoutées en vague 3)\n", count)
	return nil
}
