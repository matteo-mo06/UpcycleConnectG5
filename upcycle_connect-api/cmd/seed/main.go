package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load()

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		fmt.Println("Erreur : ADMIN_EMAIL et ADMIN_PASSWORD doivent être définis dans le .env")
		os.Exit(1)
	}

	db := openDB()
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Erreur lors du hash du mot de passe :", err)
		os.Exit(1)
	}

	// Récupère l'id du rôle admin en base (pas d'UUID hardcodé)
	var adminRoleID string
	err = db.QueryRow("SELECT id_role FROM role WHERE name_role = 'admin'").Scan(&adminRoleID)
	if err != nil {
		fmt.Println("Rôle 'admin' introuvable en base. Les migrations ont-elles été appliquées ?")
		os.Exit(1)
	}

	var userID string
	err = db.QueryRow("SELECT id_user FROM user WHERE email = ?", email).Scan(&userID)

	if err == sql.ErrNoRows {
		// L'utilisateur n'existe pas → on le crée
		userID = uuid.New().String()
		_, err = db.Exec(
			`INSERT INTO user (id_user, email, password_user, first_name, last_name, upcycling_score, premium, created_at, status)
			 VALUES (?, ?, ?, 'Admin', 'Upcycle', 0, 0, ?, 'active')`,
			userID, email, string(hash), time.Now(),
		)
		if err != nil {
			fmt.Println("Erreur lors de la création de l'admin :", err)
			os.Exit(1)
		}
		_, err = db.Exec("INSERT INTO user_role (id_user, id_role) VALUES (?, ?)", userID, adminRoleID)
		if err != nil {
			fmt.Println("Erreur lors de l'attribution du rôle admin :", err)
			os.Exit(1)
		}
		fmt.Printf("Admin créé avec succès : %s\n", email)

	} else if err == nil {
		// L'utilisateur existe déjà (ex: admin hardcodé en migration) → on met à jour le mot de passe
		_, err = db.Exec("UPDATE user SET password_user = ? WHERE id_user = ?", string(hash), userID)
		if err != nil {
			fmt.Println("Erreur lors de la mise à jour du mot de passe :", err)
			os.Exit(1)
		}
		// S'assure qu'il a bien le rôle admin
		var roleCount int
		db.QueryRow("SELECT COUNT(*) FROM user_role WHERE id_user = ? AND id_role = ?", userID, adminRoleID).Scan(&roleCount)
		if roleCount == 0 {
			db.Exec("INSERT INTO user_role (id_user, id_role) VALUES (?, ?)", userID, adminRoleID)
		}
		fmt.Printf("Mot de passe admin mis à jour pour : %s\n", email)

	} else {
		fmt.Println("Erreur lors de la vérification de l'utilisateur :", err)
		os.Exit(1)
	}

	fmt.Println("Pense à supprimer ADMIN_EMAIL et ADMIN_PASSWORD de ton .env !")
}

func openDB() *sql.DB {
	host     := os.Getenv("DB_HOST")
	port     := os.Getenv("DB_PORT")
	user     := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname   := os.Getenv("DB_NAME")

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
