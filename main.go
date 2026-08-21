package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type URL struct {
	Id        int    `json:"id"`
	Link      string `json:"link"`
	ShortLink string `json:"shortLink"`
}

var links []URL

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{shortLink}", ListingHandler)
	mux.HandleFunc("POST /shortener", CreatingHandler)

	fmt.Println("Sunucu başlatılıyor")
	http.ListenAndServe(":8080", mux)
}

func CreatingHandler(w http.ResponseWriter, r *http.Request) {
	var link URL
	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		http.Error(w, "Geçersiz JSON formatı: ", http.StatusBadRequest)
		return
	}

	fmt.Println("Link başarıyla alındı!")
	links = append(links, link)
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
