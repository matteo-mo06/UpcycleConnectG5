package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v7"
)

func seedAdvice(db *sql.DB, s *state, r *refs) error {
	const total = 35
	for i := 0; i < total; i++ {
		th := themes[i%len(themes)]
		id := newID()
		creator := pick(s.userIDs)
		createdAt := gofakeit.DateRange(daysAgo(300), daysAgo(1))
		dateAdvice := gofakeit.DateRange(daysAgo(300), daysAgo(1))
		description := fmt.Sprintf("%s : astuce pratique liée à %s.", th.Advice, th.Category)
		deletedAt := ""
		if percent(5) {
			deletedAt = sqlDateTime(gofakeit.DateRange(daysAgo(30), daysAgo(1)))
		}
		_, err := db.Exec(`
			INSERT INTO advice (id_advice, title, date_advice, description, id_creator, created_at, id_category, priority, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, th.Advice, sqlDate(dateAdvice), description, creator, sqlDateTime(createdAt),
			r.categoryIDByName[th.Category], randInt(1, 3), nullStr(deletedAt),
		)
		if err != nil {
			return err
		}
		s.adviceIDs = append(s.adviceIDs, id)
	}
	fmt.Printf("advice: %d créés\n", total)
	return nil
}

var projectStatuses = []string{"open", "open", "open", "in_progress", "completed", "pending"}

func seedProjects(db *sql.DB, s *state) error {
	const total = 45
	for i := 0; i < total; i++ {
		th := themes[i%len(themes)]
		id := newID()
		creator := pick(s.userIDs)
		title := fmt.Sprintf("%s - %s", th.Project, gofakeit.City())
		description := fmt.Sprintf("%s, %s.", th.Project, pick(announcementDescriptionSuffixes))
		start := gofakeit.DateRange(daysAgo(200), daysFromNow(30))
		end := start.AddDate(0, 0, randInt(30, 180))
		capacity := randInt(4, 20)
		createdAt := gofakeit.DateRange(daysAgo(200), daysAgo(1))

		status := pick(projectStatuses)
		rejectionReason := ""
		if percent(10) {
			status = "rejected"
			rejectionReason = "Projet non conforme aux règles de la plateforme"
		}
		isDeleted := percent(5)
		deletedAt := ""
		if isDeleted {
			deletedAt = sqlDateTime(gofakeit.DateRange(daysAgo(30), daysAgo(1)))
		}

		_, err := db.Exec(`
			INSERT INTO project (id_project, start_date_project, end_date, capacity, id_creator,
				title_project, description_project, location_project, status, rejection_reason, created_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, sqlDate(start), sqlDate(end), capacity, creator,
			title, description, gofakeit.City(), status, nullStr(rejectionReason), sqlDateTime(createdAt), nullStr(deletedAt),
		)
		if err != nil {
			return err
		}
		s.projectIDs = append(s.projectIDs, id)
		s.projectCreator[id] = creator
		s.projectDeleted[id] = isDeleted
	}
	fmt.Printf("project: %d créés\n", total)
	return nil
}

var formationLevels = []string{"beginner", "intermediate", "advanced"}

func seedFormations(db *sql.DB, s *state, r *refs) error {
	const total = 45
	teachers := append(append([]string{}, s.artisanIDs...), s.salarieIDs...)
	if len(teachers) == 0 {
		teachers = s.userIDs
	}

	for i := 0; i < total; i++ {
		th := themes[i%len(themes)]
		id := newID()
		creator := pick(s.userIDs)
		formateur := pick(teachers)
		title := fmt.Sprintf("%s - %s", th.Formation, gofakeit.City())
		description := fmt.Sprintf("%s, %s.", th.Formation, pick(announcementDescriptionSuffixes))
		dateFormation := gofakeit.DateRange(daysAgo(60), daysFromNow(120))
		capacity := randInt(6, 30)
		duration := randInt(2, 40)
		price := randMoney(20, 100)
		createdAt := gofakeit.DateRange(daysAgo(200), daysAgo(1))

		status := "approved"
		rejectionReason := ""
		if percent(10) {
			status = "rejected"
			rejectionReason = "Contenu pédagogique insuffisant, merci de compléter votre dossier"
		} else if percent(5) {
			status = "pending"
		}
		isDeleted := percent(5)
		deletedAt := ""
		if isDeleted {
			deletedAt = sqlDateTime(gofakeit.DateRange(daysAgo(30), daysAgo(1)))
		}

		_, err := db.Exec(`
			INSERT INTO formation (id_formation, id_creator, id_formateur, title_formation, description_formation,
				date_formation, capacity, level, duration_hours, status, rejection_reason, id_category, created_at,
				deleted_at, price, syllabus, objectives, address_formation, city_formation, postal_formation)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, creator, formateur, title, description,
			sqlDateTime(dateFormation), capacity, pick(formationLevels), duration, status, nullStr(rejectionReason),
			r.categoryIDByName[th.Category], sqlDateTime(createdAt),
			nullStr(deletedAt), price, "Programme : "+description, "Être capable de "+th.Formation,
			gofakeit.StreetName(), gofakeit.City(), gofakeit.Zip(),
		)
		if err != nil {
			return err
		}
		s.formationIDs = append(s.formationIDs, id)
		s.formationCreator[id] = creator
		s.formationDeleted[id] = isDeleted
	}
	fmt.Printf("formation: %d créés\n", total)
	return nil
}

func seedEvents(db *sql.DB, s *state) error {
	const total = 35
	prices := []int{0, 500, 1000, 1500, 2000, 3000, 5000}
	for i := 0; i < total; i++ {
		id := newID()
		creator := pick(s.userIDs)
		title := fmt.Sprintf("%s - %s", pick(eventTitles), gofakeit.City())
		description := fmt.Sprintf("%s organisé à %s.", pick(eventTitles), gofakeit.City())
		dateEvent := gofakeit.DateRange(daysAgo(30), daysFromNow(120))
		capacity := randInt(10, 100)

		status := "approved"
		rejectionReason := ""
		if percent(10) {
			status = "rejected"
			rejectionReason = "Informations insuffisantes sur le lieu de l'événement"
		}

		_, err := db.Exec(`
			INSERT INTO event (id_event, title_event, description_event, date_event, capacity, price_cents,
				id_creator, status, rejection_reason, address_event, city_event, postal_event)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, title, description, sqlDateTime(dateEvent), capacity, pick(prices),
			creator, status, nullStr(rejectionReason),
			gofakeit.StreetName(), gofakeit.City(), gofakeit.Zip(),
		)
		if err != nil {
			return err
		}
		s.eventIDs = append(s.eventIDs, id)
		s.eventCreator[id] = creator
	}
	fmt.Printf("event: %d créés\n", total)
	return nil
}

func seedAnnouncements(db *sql.DB, s *state, r *refs) error {
	const total = 280
	const vendu = 45
	const refusee = 15
	const supprimee = 15
	const enAttente = 28

	for i := 0; i < total; i++ {
		th := themes[i%len(themes)]
		id := newID()
		owner := pick(s.userIDs)
		baseTitle := pick(th.Announcements)
		title := fmt.Sprintf("%s - %s", baseTitle, gofakeit.City())
		description := fmt.Sprintf("%s, %s.", baseTitle, pick(announcementDescriptionSuffixes))
		availability := gofakeit.DateRange(daysAgo(60), daysFromNow(90))
		createdAt := gofakeit.DateRange(daysAgo(180), daysAgo(1))

		typeAnnouncement := "don"
		price := 0.0
		if percent(30) {
			typeAnnouncement = "vente"
			price = randMoney(0.01, 150)
		}

		request := 0
		if percent(10) {
			request = 1
		}

		var state_, rejectionReason, deletedAt string
		var buyer string
		switch {
		case i < vendu:
			state_ = "Vendu"
			buyer = pickExcept(s.userIDs, owner)
		case i < vendu+refusee:
			state_ = "Refusée"
			rejectionReason = "Description insuffisante ou photo manquante"
		case i < vendu+refusee+supprimee:
			state_ = "Supprimée"
			deletedAt = sqlDateTime(gofakeit.DateRange(daysAgo(30), daysAgo(1)))
		case i < vendu+refusee+supprimee+enAttente:
			state_ = "En attente"
		default:
			state_ = "Active"
		}

		_, err := db.Exec(`
			INSERT INTO announcement (id_announcement, address_annoucement, city, postal, description_annoucement,
				availability_date, price, request, state_annoucement, id_category, title_announcement,
				type_announcement, condition_announcement, id_buyer, deleted_at, created_at, rejection_reason)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, gofakeit.StreetName(), gofakeit.City(), gofakeit.Zip(), description,
			sqlDate(availability), price, request, state_, r.categoryIDByName[th.Category], title,
			typeAnnouncement, pick(conditionValues), nullStr(buyer), nullStr(deletedAt), sqlDateTime(createdAt), nullStr(rejectionReason),
		)
		if err != nil {
			return err
		}
		s.announcementIDs = append(s.announcementIDs, id)
		s.announcementOwner[id] = owner
		s.announcementCategory[id] = r.categoryIDByName[th.Category]
		s.announcementState[id] = state_
		s.announcementDeleted[id] = deletedAt != ""
		if buyer != "" {
			s.announcementBuyer[id] = buyer
		}
	}
	fmt.Printf("announcement: %d créées (%d Vendu, %d Refusée, %d Supprimée, %d En attente)\n", total, vendu, refusee, supprimee, enAttente)
	return nil
}

func seedProfessionalRequests(db *sql.DB, s *state) error {
	total := 30
	if len(s.regularIDs) < total {
		total = len(s.regularIDs)
	}
	rejected := total * 30 / 100
	pending := total * 20 / 100

	for i := 0; i < total; i++ {
		id := newID()
		userID := s.regularIDs[i]
		createdAt := gofakeit.DateRange(daysAgo(120), daysAgo(1))

		status := "approved"
		if i < rejected {
			status = "rejected"
		} else if i < rejected+pending {
			status = "pending"
		}

		processedAt := ""
		if status != "pending" {
			processedAt = sqlDateTime(createdAt.AddDate(0, 0, randInt(1, 10)))
		}

		_, err := db.Exec(`
			INSERT INTO professional_request (id_request, id_user, status, created_at, processed_at)
			VALUES (?, ?, ?, ?, ?)`,
			id, userID, status, sqlDateTime(createdAt), nullStr(processedAt),
		)
		if err != nil {
			return err
		}
		s.professionalRequestIDs = append(s.professionalRequestIDs, id)
	}
	fmt.Printf("professional_request: %d créées\n", total)
	return nil
}

func seedDocuments(db *sql.DB, s *state) error {
	const total = 30
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	type parent struct {
		id    string
		owner string
	}
	var pool []parent
	for _, id := range s.announcementIDs {
		pool = append(pool, parent{id, s.announcementOwner[id]})
	}
	for _, id := range s.projectIDs {
		pool = append(pool, parent{id, s.projectCreator[id]})
	}
	for _, id := range s.formationIDs {
		pool = append(pool, parent{id, s.formationCreator[id]})
	}

	for i := 0; i < total; i++ {
		p := pick(pool)
		id := newID()
		link := fmt.Sprintf("%s/uploads/demo-%s.jpg", baseURL, newID())
		_, err := db.Exec(
			"INSERT INTO document (id_document, id_user, category, link) VALUES (?, ?, ?, ?)",
			id, p.owner, p.id, link,
		)
		if err != nil {
			return err
		}
	}
	fmt.Printf("document: %d créés\n", total)
	return nil
}

func seedTopics(db *sql.DB, s *state, r *refs) error {
	const total = 70
	for i := 0; i < total; i++ {
		id := newID()
		author := pick(s.userIDs)
		title := fmt.Sprintf("%s (%d)", pick(topicTitles), i+1)
		createdAt := gofakeit.DateRange(daysAgo(200), daysAgo(1))

		isDeleted := percent(5)
		deletedAt := ""
		if isDeleted {
			deletedAt = sqlDateTime(gofakeit.DateRange(daysAgo(20), daysAgo(1)))
		}
		moderatedBy := ""
		moderationMessage := ""
		if !isDeleted && percent(10) {
			moderatedBy = r.adminUserID
			moderationMessage = "Sujet recadré, merci de rester courtois"
		}

		_, err := db.Exec(`
			INSERT INTO topic (id_topic, deleted_at, title_topic, id_author, created_at, moderated_by, moderation_message)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			id, nullStr(deletedAt), title, author, sqlDateTime(createdAt), nullStr(moderatedBy), nullStr(moderationMessage),
		)
		if err != nil {
			return err
		}
		s.topicIDs = append(s.topicIDs, id)
		s.topicAuthor[id] = author
		s.topicDeleted[id] = isDeleted
	}
	fmt.Printf("topic: %d créés\n", total)
	return nil
}

func seedPosts(db *sql.DB, s *state, r *refs) error {
	const total = 280
	postsByTopic := make(map[string][]string)

	for i := 0; i < total; i++ {
		id := newID()
		author := pick(s.userIDs)
		topicID := pick(s.topicIDs)
		body := pick(postBodies)
		createdAt := gofakeit.DateRange(daysAgo(180), daysAgo(1))

		parentPost := ""
		if existing := postsByTopic[topicID]; len(existing) > 0 && percent(40) {
			parentPost = pick(existing)
		}

		isDeleted := percent(5)
		deletedAt := ""
		if isDeleted {
			deletedAt = sqlDateTime(gofakeit.DateRange(daysAgo(20), daysAgo(1)))
		}
		moderatedBy := ""
		moderationMessage := ""
		if !isDeleted && percent(5) {
			moderatedBy = r.adminUserID
			moderationMessage = "Message modéré, propos inappropriés"
		}

		_, err := db.Exec(`
			INSERT INTO post (id_post, deleted_at, body_post, id_author, id_parent_post, created_at, moderated_by, moderation_message)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			id, nullStr(deletedAt), body, author, nullStr(parentPost), sqlDateTime(createdAt), nullStr(moderatedBy), nullStr(moderationMessage),
		)
		if err != nil {
			return err
		}
		s.postIDs = append(s.postIDs, id)
		s.postTopic[id] = topicID
		s.postAuthor[id] = author
		s.postDeleted[id] = isDeleted
		postsByTopic[topicID] = append(postsByTopic[topicID], id)
	}
	fmt.Printf("post: %d créés\n", total)
	return nil
}

var advertisementStates = []string{"pending", "approved", "rejected", "active", "expired"}

func seedAdvertisements(db *sql.DB, s *state, r *refs) error {
	const perState = 14
	total := 0
	advertisers := append(append([]string{}, s.artisanIDs...), s.salarieIDs...)
	if len(advertisers) == 0 {
		advertisers = s.userIDs
	}

	for _, state_ := range advertisementStates {
		for i := 0; i < perState; i++ {
			id := newID()
			owner := pick(advertisers)
			if percent(30) {
				owner = pick(s.userIDs)
			}
			plan := pick(r.advertisementPlans)
			title := fmt.Sprintf("Publicité %s", pick(themes).Category)
			createdAt := gofakeit.DateRange(daysAgo(150), daysAgo(1))

			var approvedAt, paidAt, expiresAt, rejectionReason string
			switch state_ {
			case "rejected":
				rejectionReason = "Visuel non conforme à la charte de la plateforme"
			case "approved":
				approvedAt = sqlDateTime(createdAt.AddDate(0, 0, randInt(1, 3)))
			case "active":
				approvedAt = sqlDateTime(createdAt.AddDate(0, 0, randInt(1, 3)))
				paidAt = sqlDateTime(createdAt.AddDate(0, 0, randInt(3, 5)))
				expiresAt = sqlDate(daysFromNow(randInt(3, plan.Weeks*7)))
			case "expired":
				approvedAt = sqlDateTime(createdAt.AddDate(0, 0, randInt(1, 3)))
				paidAt = sqlDateTime(createdAt.AddDate(0, 0, randInt(3, 5)))
				expiresAt = sqlDate(daysAgo(randInt(1, 30)))
			}

			_, err := db.Exec(`
				INSERT INTO advertisement (id_advertisement, id_user, title, image_url, link_url, state, price_cents,
					rejection_reason, created_at, approved_at, paid_at, expires_at, plan_id)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				id, owner, title, fmt.Sprintf("https://picsum.photos/seed/%s/600/300", id[:8]), nil, state_, plan.PriceCents,
				nullStr(rejectionReason), sqlDateTime(createdAt), nullStr(approvedAt), nullStr(paidAt), nullStr(expiresAt), plan.ID,
			)
			if err != nil {
				return err
			}
			s.advertisementIDs = append(s.advertisementIDs, id)
			total++
		}
	}
	fmt.Printf("advertisement: %d créées (%d par état sur %d états)\n", total, perState, len(advertisementStates))
	return nil
}

func seedUserNotifications(db *sql.DB, s *state) error {
	const total = 280
	seen := newPairSet()
	count := 0
	attempts := 0
	for count < total && attempts < total*20 {
		attempts++
		userID := pick(s.userIDs)
		notifID := pick(s.notificationIDs)
		if !seen.add(userID, notifID) {
			continue
		}
		if _, err := db.Exec(
			"INSERT INTO user_notification (id_user, id_notification) VALUES (?, ?)", userID, notifID,
		); err != nil {
			return err
		}
		count++
	}
	fmt.Printf("user_notification: %d créés\n", count)
	return nil
}

func seedSubscriptions(db *sql.DB, s *state) error {
	const total = 35
	for i := 0; i < total; i++ {
		id := newID()
		owner := pick(s.userIDs)
		start := gofakeit.DateRange(daysAgo(300), daysAgo(10))
		end := start.AddDate(1, 0, 0)
		cancelled := 0
		cancelledAt := ""
		if percent(15) {
			cancelled = 1
			cancelledAt = sqlDateTime(start.AddDate(0, randInt(1, 6), 0))
		}

		_, err := db.Exec(`
			INSERT INTO subscription (id_subscription, start_timestamp, end_timestamp, cancelled, cancelled_at,
				stripe_subscription_id, stripe_customer_id)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			id, sqlDateTime(start), sqlDateTime(end), cancelled, nullStr(cancelledAt),
			fmt.Sprintf("sub_demo_%s", id[:12]), fmt.Sprintf("cus_demo_%s", id[:12]),
		)
		if err != nil {
			return err
		}
		s.subscriptionIDs = append(s.subscriptionIDs, id)
		s.subscriptionUser[id] = owner
	}
	fmt.Printf("subscription: %d créés\n", total)
	return nil
}

func seedPayements(db *sql.DB, s *state, r *refs) error {
	insert := func(kind, announcementID, buyerID, formationID, subscriptionID string, amountCents int, commissionCents int, paidAt string) (string, error) {
		id := newID()
		_, err := db.Exec(`
			INSERT INTO payement (id_payement, amount_cents, currency, provider_payement, provider_ref, status_payement,
				created_at, paid_at, announcement_id, buyer_id, commission_amount_cents, formation_id, subscription_id)
			VALUES (?, ?, 'eur', 'stripe', ?, 'paid', ?, ?, ?, ?, ?, ?, ?)`,
			id, amountCents, fmt.Sprintf("pi_demo_%s", id[:12]), paidAt, paidAt,
			nullStr(announcementID), nullStr(buyerID), commissionCents, nullStr(formationID), nullStr(subscriptionID),
		)
		return id, err
	}

	for _, annID := range s.announcementIDs {
		buyerID, ok := s.announcementBuyer[annID]
		if !ok {
			continue
		}
		amount := int(randMoney(5, 150) * 100)
		commission := amount * int(r.commissionRate) / 100
		paidAt := sqlDateTime(gofakeit.DateRange(daysAgo(60), daysAgo(1)))
		id, err := insert("announcement", annID, buyerID, "", "", amount, commission, paidAt)
		if err != nil {
			return err
		}
		s.payements.announcementPay = append(s.payements.announcementPay, id)
		s.payements.all = append(s.payements.all, id)
		s.payements.buyerOf[id] = buyerID
	}

	for _, subID := range s.subscriptionIDs {
		userID := s.subscriptionUser[subID]
		amount := 1500
		paidAt := sqlDateTime(gofakeit.DateRange(daysAgo(90), daysAgo(1)))
		id, err := insert("subscription", "", userID, "", subID, amount, 0, paidAt)
		if err != nil {
			return err
		}
		s.payements.subscriptionPay = append(s.payements.subscriptionPay, id)
		s.payements.all = append(s.payements.all, id)
		s.payements.buyerOf[id] = userID
	}

	const formationPayCount = 45
	for i := 0; i < formationPayCount; i++ {
		formationID := pick(s.formationIDs)
		buyerID := pick(s.userIDs)
		amount := int(randMoney(20, 100) * 100)
		paidAt := sqlDateTime(gofakeit.DateRange(daysAgo(90), daysAgo(1)))
		id, err := insert("formation", "", buyerID, formationID, "", amount, 0, paidAt)
		if err != nil {
			return err
		}
		s.payements.formationPay = append(s.payements.formationPay, id)
		s.payements.all = append(s.payements.all, id)
		s.payements.buyerOf[id] = buyerID
		s.payements.formationOf[id] = formationID
	}

	const eventPayCount = 35
	for i := 0; i < eventPayCount; i++ {
		buyerID := pick(s.userIDs)
		amount := pick([]int{500, 1000, 1500, 2000, 3000})
		paidAt := sqlDateTime(gofakeit.DateRange(daysAgo(90), daysAgo(1)))
		id, err := insert("event", "", buyerID, "", "", amount, 0, paidAt)
		if err != nil {
			return err
		}
		s.payements.eventPay = append(s.payements.eventPay, id)
		s.payements.all = append(s.payements.all, id)
		s.payements.buyerOf[id] = buyerID
	}

	const projectPayCount = 30
	for i := 0; i < projectPayCount; i++ {
		buyerID := pick(s.userIDs)
		amount := randInt(10, 80) * 100
		paidAt := sqlDateTime(gofakeit.DateRange(daysAgo(90), daysAgo(1)))
		id, err := insert("project", "", buyerID, "", "", amount, 0, paidAt)
		if err != nil {
			return err
		}
		s.payements.projectPay = append(s.payements.projectPay, id)
		s.payements.all = append(s.payements.all, id)
		s.payements.buyerOf[id] = buyerID
	}

	total := len(s.payements.announcementPay) + len(s.payements.subscriptionPay) + len(s.payements.formationPay) + len(s.payements.eventPay) + len(s.payements.projectPay)
	fmt.Printf("payement: %d créés\n", total)
	return nil
}
