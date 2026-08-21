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

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprint(os.Stderr, "Veritabanına bağlanılamadı")
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
	links = append(links, link)
	fmt.Println("Link başarıyla oluşturuldu!")

	json.NewEncoder(w).Encode(link)
}

func ListingHandler(w http.ResponseWriter, r *http.Request) {
	shortLink := r.PathValue("shortLink")
	for i := range links {
		if links[i].ShortLink == shortLink {
			http.Redirect(w, r, links[i].Link, 302)
		}
	}
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
