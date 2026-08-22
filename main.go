package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/yigitarac/go-url-shortener/handlers"
	"github.com/yigitarac/go-url-shortener/middleware"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env dosyası yüklenemedi")
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /p/{shortLink}", middleware.FirstMiddleware(handlers.ListingHandler))
	mux.HandleFunc("POST /shortener", middleware.FirstMiddleware(handlers.CreatingHandler))

	fileServer := http.FileServer(http.Dir("./ui"))
	mux.Handle("/", fileServer)

	fmt.Println("Sunucu başlatılıyor")

	databaseURL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprint(os.Stderr, "Veritabanına bağlanılamadı")
		return
	}

	defer conn.Close(context.Background())
	handlers.Conn = conn
	http.ListenAndServe(":8080", mux)
}
