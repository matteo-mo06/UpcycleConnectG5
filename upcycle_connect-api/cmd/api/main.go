package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ValidateEnv()

	config.Conn = config.NewDB()

	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "./uploads"
	}
	os.MkdirAll(uploadsDir, 0755)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	os.MkdirAll(config.InvoicesDir(), 0755)

	router.InitRoutes()

	go func() {
		for {
			if err := db.ExpireOverdueAdvertisements(); err != nil {
				fmt.Println("ExpireOverdueAdvertisements error:", err)
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	port := os.Getenv("PORT")
	fmt.Println("Listening at http://localhost:" + port)
	http.ListenAndServe(":"+port, middleware.CORS(http.DefaultServeMux))
}
