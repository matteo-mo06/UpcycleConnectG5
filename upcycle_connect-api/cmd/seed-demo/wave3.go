package main

import (
	"database/sql"
	"fmt"
)

type sanctionSpec struct {
	sanctionType string
	durationDays int
}

func seedUserSanctions(db *sql.DB, s *state, r *refs, points map[string]int) error {
	const warningCount = 16
	const suspensionCount = 10
	const banCount = 4
	suspensionDurations := []int{7, 14, 14, 30, 30}

	var specs []sanctionSpec
	for i := 0; i < warningCount; i++ {
		specs = append(specs, sanctionSpec{"warning", 0})
	}
	for i := 0; i < suspensionCount; i++ {
		specs = append(specs, sanctionSpec{"suspension", suspensionDurations[i%len(suspensionDurations)]})
	}
	for i := 0; i < banCount; i++ {
		specs = append(specs, sanctionSpec{"ban", 0})
	}

	if len(s.reportResolved) < len(specs) {
		return fmt.Errorf("pas assez de signalements résolus (%d) pour créer %d sanctions", len(s.reportResolved), len(specs))
	}

	for i, spec := range specs {
		reportID := s.reportResolved[i]
		targetUser := s.reportedUserFor[reportID]
		id := newID()
		createdAt := daysAgo(randInt(1, 20))

		expiresAt := ""
		var durationDays any
		if spec.durationDays > 0 {
			expiresAt = sqlDateTime(createdAt.AddDate(0, 0, spec.durationDays))
			durationDays = spec.durationDays
		}

		_, err := db.Exec(`
			INSERT INTO user_sanctions (id_sanction, id_user, id_admin, id_report, type, reason, duration_days, expires_at, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, targetUser, r.adminUserID, reportID, spec.sanctionType, pick(sanctionReasons),
			durationDays, nullStr(expiresAt), sqlDateTime(createdAt),
		)
		if err != nil {
			return err
		}
		s.userSanctionIDs = append(s.userSanctionIDs, id)

		switch spec.sanctionType {
		case "warning":
			if _, err := db.Exec(
				"INSERT INTO score_log (id, id_user, action_type, id_reference, points) VALUES (?, ?, 'sanction_warning', ?, ?)",
				newID(), targetUser, id, points["sanction_warning"],
			); err != nil {
				return err
			}
		case "suspension":
			if _, err := db.Exec(
				"INSERT INTO score_log (id, id_user, action_type, id_reference, points) VALUES (?, ?, 'sanction_suspension', ?, ?)",
				newID(), targetUser, id, points["sanction_suspension"],
			); err != nil {
				return err
			}
			if _, err := db.Exec("UPDATE user SET status = 'suspended' WHERE id_user = ?", targetUser); err != nil {
				return err
			}
		case "ban":
			if _, err := db.Exec("UPDATE user SET status = 'blacklisted' WHERE id_user = ?", targetUser); err != nil {
				return err
			}
		}
	}

	fmt.Printf("user_sanctions: %d créés (%d warning, %d suspension, %d ban) + score_log associés\n", len(specs), warningCount, suspensionCount, banCount)
	return nil
}

func seedPayementJunctions(db *sql.DB, s *state) error {
	for _, paymentID := range s.payements.eventPay {
		eventID := pick(s.eventIDs)
		if _, err := db.Exec(
			"INSERT INTO payement_event (id_payement, id_event) VALUES (?, ?)", paymentID, eventID,
		); err != nil {
			return err
		}
	}
	fmt.Printf("payement_event: %d créés\n", len(s.payements.eventPay))

	for _, paymentID := range s.payements.formationPay {
		formationID := s.payements.formationOf[paymentID]
		if _, err := db.Exec(
			"INSERT INTO payement_formation (id_payement, id_formation) VALUES (?, ?)", paymentID, formationID,
		); err != nil {
			return err
		}
	}
	fmt.Printf("payement_formation: %d créés\n", len(s.payements.formationPay))

	for _, paymentID := range s.payements.projectPay {
		projectID := pick(s.projectIDs)
		if _, err := db.Exec(
			"INSERT INTO payement_project (id_payement, id_project) VALUES (?, ?)", paymentID, projectID,
		); err != nil {
			return err
		}
	}
	fmt.Printf("payement_project: %d créés\n", len(s.payements.projectPay))

	for _, paymentID := range s.payements.all {
		buyerID := s.payements.buyerOf[paymentID]
		if _, err := db.Exec(
			"INSERT INTO user_payement (id_user, id_payement) VALUES (?, ?)", buyerID, paymentID,
		); err != nil {
			return err
		}
	}
	fmt.Printf("user_payement: %d créés\n", len(s.payements.all))

	for _, paymentID := range s.payements.formationPay {
		buyerID := s.payements.buyerOf[paymentID]
		formationID := s.payements.formationOf[paymentID]
		if _, err := db.Exec(
			"INSERT INTO user_formation_inscription_payement (id_user, id_formation, id_payement) VALUES (?, ?, ?)",
			buyerID, formationID, paymentID,
		); err != nil {
			return err
		}
	}
	fmt.Printf("user_formation_inscription_payement: %d créés\n", len(s.payements.formationPay))

	return nil
}

func seedUserSubscriptionAndPlan(db *sql.DB, s *state, r *refs) error {
	for _, subID := range s.subscriptionIDs {
		userID := s.subscriptionUser[subID]
		if _, err := db.Exec(
			"INSERT INTO user_subscription (id_user, id_subscription) VALUES (?, ?)", userID, subID,
		); err != nil {
			return err
		}
		if _, err := db.Exec(
			"INSERT INTO sub_sub_plan (id_subscription, id_plan) VALUES (?, ?)", subID, r.subscriptionPlanID,
		); err != nil {
			return err
		}
	}
	fmt.Printf("user_subscription: %d créés, sub_sub_plan: %d créés\n", len(s.subscriptionIDs), len(s.subscriptionIDs))
	return nil
}
