package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type URL struct {
	Id        int    `json:"id"`
	Link      string `json:"link"`
	ShortLink string `json:"shortLink"`
}

var links []URL
var conn *pgx.Conn

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env dosyası yüklenemedi")
		return
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{shortLink}", ListingHandler)
	mux.HandleFunc("POST /shortener", CreatingHandler)

	fmt.Println("Sunucu başlatılıyor")

	databaseURL := os.Getenv("DATABASE_URL")

	conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprint(os.Stderr, "Veritabanına bağlanılamadı")
		return
	}

	defer conn.Close(context.Background())

	http.ListenAndServe(":8080", mux)

}

func CreatingHandler(w http.ResponseWriter, r *http.Request) {
	var link URL
	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		http.Error(w, "Geçersiz JSON formatı: ", http.StatusBadRequest)
		return
	}

	link.ShortLink = createShortLink()
	fmt.Println("Link başarıyla oluşturuldu!")

	_, err = conn.Exec(context.Background(), "INSERT INTO links(link, shortLink) VALUES ($1, $2)", link.Link, link.ShortLink)
	if err != nil {
		fmt.Println("Kısaltılan link veritabanına eklenemedi")
		return
	}

	json.NewEncoder(w).Encode(link)
}

func ListingHandler(w http.ResponseWriter, r *http.Request) {
	shortLink := r.PathValue("shortLink")
	var longLink string
	err := conn.QueryRow(context.Background(), "SELECT link FROM links WHERE shortLink = $1", shortLink).Scan(&longLink)
	if err != nil {
		fmt.Println("Uzun link bulunamadı")
		return
	}
	http.Redirect(w, r, longLink, 302)
}

func createShortLink() string {
	var shortLink []byte
	characters := "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuopasdfghjklizxcvbnm123456789"
	for i := 0; i < 6; i++ {
		index := rand.Intn(len(characters))
		shortLink = append(shortLink, characters[index])
	}
	return string(shortLink)
}
