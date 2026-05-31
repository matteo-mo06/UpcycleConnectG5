package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var JWTSecret = os.Getenv("JWT_SECRET")

const RoleAdmin   = "admin"
const RoleArtisan = "artisan"
const RoleSalarie = "salarie"
const RoleUser    = "user"

const AdminRoleName   = RoleAdmin
const ArtisanRoleName = RoleArtisan

var Conn *sql.DB

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
