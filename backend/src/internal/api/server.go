package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Store struct {
	Directory string
	Stop      chan struct{}
}

func (store *Store) Serve(port int) {
	handler := store.NewHandler()

	addr := fmt.Sprintf(":%d", port)

	server := &http.Server{Addr: addr, Handler: handler}

	go func() {
		log.Printf("MAPATLAPI serving HTTP on %s\n", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("MAPATLAPI server stopped unexpectedly: %v", err)
		}
	}()

	<-store.Stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("MAPATLAPI shutdown failed: %v\n", err)
	}
}
