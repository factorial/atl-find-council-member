package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Store struct {
	Directory string
	Stop      chan struct{}
}
type Interface interface{}

func (store *Store) Serve(port int) {
	r := chi.NewRouter()

	r.Use(newCorsHandler())

	r.Post("/address", NewSearchAddressHandler(store))
	r.Post("/record", NewRecordRequestHandler(store))

	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: addr, Handler: r}

	go func() {
		log.Printf("—APISERVICE— serving HTTP on %s\n", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("API Server stopped unexpectedly: %v", err)
		}
	}()

	<-store.Stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("—APISERVICE— shutdown failed: %v\n", err)
	}
}

func newCorsHandler() func(http.Handler) http.Handler {
	options := cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	cors := cors.New(options)
	return cors.Handler
}
