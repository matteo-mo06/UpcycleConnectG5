package main

import (
	"fmt"
	"net/http"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.Conn = config.NewDB()

	router.InitRoutes()

	fmt.Println("Listening at http://localhost:8084")
	http.ListenAndServe(":8084", nil)
}
