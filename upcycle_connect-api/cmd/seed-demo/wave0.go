package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/crypto/bcrypt"
)

const demoPassword = "Demo1234!"

func seedUsers(db *sql.DB, s *state, r *refs) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	const total = 60
	const artisanCount = 16
	const salarieCount = 8

	for i := 0; i < total; i++ {
		id := newID()
		first := gofakeit.FirstName()
		last := gofakeit.LastName()
		email := fmt.Sprintf("%s.%s%d@example.com", strings.ToLower(first), strings.ToLower(last), i)
		createdAt := gofakeit.DateRange(daysAgo(365), daysAgo(1))

		roleName := "user"
		if i < artisanCount {
			roleName = "artisan"
		} else if i < artisanCount+salarieCount {
			roleName = "salarie"
		}

		_, err := db.Exec(`
			INSERT INTO user (id_user, email, password_user, first_name, last_name, upcycling_score, premium, created_at, status, email_verified)
			VALUES (?, ?, ?, ?, ?, 0, 0, ?, 'active', 1)`,
			id, email, string(hash), first, last, sqlDateTime(createdAt),
		)
		if err != nil {
			return err
		}

		roleID := r.roleIDByName[roleName]
		if _, err := db.Exec(
			"INSERT INTO user_role (id_user, id_role) VALUES (?, ?)", id, roleID,
		); err != nil {
			return err
		}

		s.userIDs = append(s.userIDs, id)
		s.userRoleName[id] = roleName
		switch roleName {
		case "artisan":
			s.artisanIDs = append(s.artisanIDs, id)
		case "salarie":
			s.salarieIDs = append(s.salarieIDs, id)
		default:
			s.regularIDs = append(s.regularIDs, id)
		}
	}

	fmt.Printf("users: %d créés (%d artisan, %d salarie, %d user)\n", total, len(s.artisanIDs), len(s.salarieIDs), len(s.regularIDs))
	return nil
}

func seedTickets(db *sql.DB) error {
	const total = 100
	for i := 0; i < total; i++ {
		ts := gofakeit.DateRange(daysAgo(180), daysFromNow(60))
		if _, err := db.Exec(
			"INSERT INTO ticket (id_ticket, timestamp_ticket) VALUES (?, ?)",
			newID(), sqlDateTime(ts),
		); err != nil {
			return err
		}
	}
	fmt.Printf("ticket: %d créés\n", total)
	return nil
}

func seedNotifications(db *sql.DB, s *state) error {
	const total = 180
	for i := 0; i < total; i++ {
		id := newID()
		view := 0
		if percent(50) {
			view = 1
		}
		if _, err := db.Exec(
			"INSERT INTO notification (id_notification, view) VALUES (?, ?)",
			id, view,
		); err != nil {
			return err
		}
		s.notificationIDs = append(s.notificationIDs, id)
	}
	fmt.Printf("notification: %d créées\n", total)
	return nil
}
