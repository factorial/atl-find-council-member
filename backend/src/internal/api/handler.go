package api

import (
	"backend/internal/api/models/address"
	"backend/internal/api/models/council"
	"backend/internal/api/models/record"
	"net/http"
	"path"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func (store *Store) NewHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(newCorsHandler())

	r.Post("/council", council.NewHandler(path.Join(store.Directory, "citycouncil.json")))
	r.Post("/address", address.NewHandler())
	r.Post("/record", record.NewHandler())

	return r
}

func newCorsHandler() func(http.Handler) http.Handler {
	options := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	cors := cors.New(options)
	return cors.Handler
}
