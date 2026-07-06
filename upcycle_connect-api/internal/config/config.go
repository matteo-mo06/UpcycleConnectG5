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

func JWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func StripeSecretKey() string     { return os.Getenv("STRIPE_SECRET_KEY") }
func StripeWebhookSecret() string { return os.Getenv("STRIPE_WEBHOOK_SECRET") }
func BaseURL() string     { return os.Getenv("BASE_URL") }
func FrontendURL() string { return os.Getenv("FRONTEND_URL") }
func PlatformURL() string { return os.Getenv("PLATFORM_URL") }

func SMTPHost() string     { return os.Getenv("SMTP_HOST") }
func SMTPPort() string     { return os.Getenv("SMTP_PORT") }
func SMTPUsername() string { return os.Getenv("SMTP_USERNAME") }
func SMTPPassword() string { return os.Getenv("SMTP_PASSWORD") }
func MailFrom() string     { return os.Getenv("MAIL_FROM") }
func SMTPTLSCert() string  { return os.Getenv("SMTP_TLS_CERT") }

func InvoicesDir() string {
	if d := os.Getenv("INVOICES_DIR"); d != "" {
		return d
	}
	return "./invoices"
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
		"STRIPE_SECRET_KEY",
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
