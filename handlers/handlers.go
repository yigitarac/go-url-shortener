package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ShortenedLink struct {
	Id        int    `json:"id"`
	Link      string `json:"link"`
	ShortLink string `json:"shortLink"`
}

var links []ShortenedLink
var Conn *pgx.Conn

func CreatingHandler(w http.ResponseWriter, r *http.Request) {
	var link ShortenedLink
	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		http.Error(w, "Geçersiz JSON formatı: ", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(link.Link)
	if err != nil {
		http.Error(w, "Geçersiz link girişi", http.StatusBadRequest)
		return
	}

	link.ShortLink = createShortLink()
	var uniqueErr *pgconn.PgError
	for {
		_, err = Conn.Exec(context.Background(), "INSERT INTO links(link, shortLink) VALUES ($1, $2)", link.Link, link.ShortLink)
		if err != nil {
			if errors.As(err, &uniqueErr) {
				if uniqueErr.Code == "23505" {
					link.ShortLink = createShortLink()
				} else {
					fmt.Println("Kritik Hata!")
					return
				}
			}
		} else {
			break
		}
	}
	fmt.Println("Link başarıyla oluşturuldu!")

	json.NewEncoder(w).Encode(link)
}

func ListingHandler(w http.ResponseWriter, r *http.Request) {
	shortLink := r.PathValue("shortLink")
	var longLink string
	err := Conn.QueryRow(context.Background(), "SELECT link FROM links WHERE shortLink = $1", shortLink).Scan(&longLink)
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
