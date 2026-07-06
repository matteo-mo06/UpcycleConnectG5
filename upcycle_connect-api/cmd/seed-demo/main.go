package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	gofakeit.Seed(42)

	db := openDB()
	defer db.Close()

	r, err := loadRefs(db)
	if err != nil {
		fmt.Println("Erreur lors du chargement des données de référence :", err)
		os.Exit(1)
	}

	s := newState()

	fmt.Println("--- Préparation (utilisateurs, tickets, notifications) ---")
	if err := seedUsers(db, s, r); err != nil {
		fail("seedUsers", err)
	}
	if err := seedTickets(db); err != nil {
		fail("seedTickets", err)
	}
	if err := seedNotifications(db, s); err != nil {
		fail("seedNotifications", err)
	}

	fmt.Println("--- Vague 1 ---")
	if err := seedAdvice(db, s, r); err != nil {
		fail("seedAdvice", err)
	}
	if err := seedProjects(db, s); err != nil {
		fail("seedProjects", err)
	}
	if err := seedFormations(db, s, r); err != nil {
		fail("seedFormations", err)
	}
	if err := seedEvents(db, s); err != nil {
		fail("seedEvents", err)
	}
	if err := seedAnnouncements(db, s, r); err != nil {
		fail("seedAnnouncements", err)
	}
	if err := seedProfessionalRequests(db, s); err != nil {
		fail("seedProfessionalRequests", err)
	}
	if err := seedDocuments(db, s); err != nil {
		fail("seedDocuments", err)
	}
	if err := seedTopics(db, s, r); err != nil {
		fail("seedTopics", err)
	}
	if err := seedPosts(db, s, r); err != nil {
		fail("seedPosts", err)
	}
	if err := seedAdvertisements(db, s, r); err != nil {
		fail("seedAdvertisements", err)
	}
	if err := seedUserNotifications(db, s); err != nil {
		fail("seedUserNotifications", err)
	}
	if err := seedSubscriptions(db, s); err != nil {
		fail("seedSubscriptions", err)
	}
	if err := seedPayements(db, s, r); err != nil {
		fail("seedPayements", err)
	}

	fmt.Println("--- Vague 2 ---")
	if err := seedAnnouncementLocker(db, s, r); err != nil {
		fail("seedAnnouncementLocker", err)
	}
	projectPairs, err := seedUserProjectInscription(db, s)
	if err != nil {
		fail("seedUserProjectInscription", err)
	}
	if err := seedUserProjectUpdown(db, s); err != nil {
		fail("seedUserProjectUpdown", err)
	}
	formationPairs, err := seedUserFormationInscription(db, s)
	if err != nil {
		fail("seedUserFormationInscription", err)
	}
	eventPairs, err := seedUserEventInscription(db, s)
	if err != nil {
		fail("seedUserEventInscription", err)
	}
	if err := seedUserAnnouncement(db, s); err != nil {
		fail("seedUserAnnouncement", err)
	}
	if err := seedUserTopic(db, s); err != nil {
		fail("seedUserTopic", err)
	}
	if err := seedUserAdvice(db, s); err != nil {
		fail("seedUserAdvice", err)
	}
	if err := seedTopicPost(db, s); err != nil {
		fail("seedTopicPost", err)
	}
	if err := seedProjectMaterialAndSteps(db, s); err != nil {
		fail("seedProjectMaterialAndSteps", err)
	}
	if err := seedFormationSteps(db, s); err != nil {
		fail("seedFormationSteps", err)
	}
	if err := seedFeaturedSlotLog(db, s); err != nil {
		fail("seedFeaturedSlotLog", err)
	}
	if err := seedReports(db, s, r); err != nil {
		fail("seedReports", err)
	}
	points, err := loadScoreActionPoints(db)
	if err != nil {
		fail("loadScoreActionPoints", err)
	}
	if err := seedScoreLog(db, s, points, projectPairs, formationPairs, eventPairs); err != nil {
		fail("seedScoreLog", err)
	}

	fmt.Println("--- Vague 3 ---")
	if err := seedUserSanctions(db, s, r, points); err != nil {
		fail("seedUserSanctions", err)
	}
	if err := seedPayementJunctions(db, s); err != nil {
		fail("seedPayementJunctions", err)
	}
	if err := seedUserSubscriptionAndPlan(db, s, r); err != nil {
		fail("seedUserSubscriptionAndPlan", err)
	}

	fmt.Println("--- Finalisation ---")
	if err := finalizeUserAggregates(db); err != nil {
		fail("finalizeUserAggregates", err)
	}

	fmt.Println("Seed de démonstration terminé avec succès.")
	fmt.Printf("Mot de passe des %d comptes de démo : %s\n", len(s.userIDs), demoPassword)
}

func fail(step string, err error) {
	fmt.Printf("Erreur lors de %s : %v\n", step, err)
	os.Exit(1)
}

func finalizeUserAggregates(db *sql.DB) error {
	if _, err := db.Exec(`
		UPDATE user u SET u.upcycling_score = COALESCE(
			(SELECT SUM(sl.points) FROM score_log sl WHERE sl.id_user = u.id_user), 0)`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(`
		UPDATE user u SET u.announcement_count = (
			SELECT COUNT(*) FROM user_announcement ua WHERE ua.id_user = u.id_user)`,
	); err != nil {
		return err
	}
	return nil
}

func openDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Erreur de connexion à la base :", err)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Impossible de joindre la base de données :", err)
		os.Exit(1)
	}
	return db
}
