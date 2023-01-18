package main

import (
	"net/http"
	"zmail/apirouter"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Use(middleware.Logger)
	r.Get("/search/{word}", func(w http.ResponseWriter, r *http.Request) {
		wordParam := chi.URLParam(r, "word")
		res := apirouter.SearchInDataBase(wordParam)

		w.Write([]byte(res))
	})
	http.ListenAndServe(":3000", r)
}
