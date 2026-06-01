package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const RoleAdmin   = "admin"
const RoleArtisan = "artisan"
const RoleSalarie = "salarie"
const RoleUser    = "user"

const AdminRoleName   = RoleAdmin
const ArtisanRoleName = RoleArtisan

var Conn *sql.DB

// JWTSecret lit la variable au moment de l'appel (après godotenv.Load).
func JWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func ValidateEnv() {
	required := []string{
		"PORT",
		"JWT_SECRET",
		"CORS_ORIGIN",
		"BASE_URL",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
	}
	for _, key := range required {
		if os.Getenv(key) == "" {
			fmt.Fprintf(os.Stderr, "FATAL: variable %s manquante dans le .env\n", key)
			os.Exit(1)
		}
	}
}

func NewDB() *sql.DB {
	host     := os.Getenv("DB_HOST")
	port     := os.Getenv("DB_PORT")
	user     := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname   := os.Getenv("DB_NAME")

	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	conn, err := sql.Open("mysql", mysqlInfo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("connected to database !")
	return conn
}
