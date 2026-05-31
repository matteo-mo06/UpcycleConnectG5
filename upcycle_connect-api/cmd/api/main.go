package main

import (
	"fmt"
	"net/http"
	"os"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.Conn = config.NewDB()

	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "./uploads"
	}
	os.MkdirAll(uploadsDir, 0755)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	router.InitRoutes()

	fmt.Println("Listening at http://localhost:8084")
	http.ListenAndServe(":8084", middleware.CORS(http.DefaultServeMux))
}
