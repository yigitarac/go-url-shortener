package middleware

import (
	"fmt"
	"net/http"
)

func FirstMiddleware(realFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method + " " + r.URL.Path)
		realFunc(w, r)
	}
}
