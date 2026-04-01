package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/handler"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	// Use PostgreSQL if DB_HOST is set, otherwise fall back to in-memory store
	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		pgStore, err := store.NewPostgresStore(
			dbHost,
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "catalog"),
			getEnv("DB_PASSWORD", "catalog123"),
			getEnv("DB_NAME", "productcatalog"),
		)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		if err := pgStore.EnsureTable(); err != nil {
			pgStore.DB.Close()
			log.Fatalf("Failed to create table: %v", err)
		}
		defer pgStore.DB.Close()

		h := handler.NewPostgresHandler(pgStore)
		h.RegisterRoutes(r)
		fmt.Printf("Product Catalog API (PostgreSQL) listening on :%s\n", port)
	} else {
		memStore := store.NewMemoryStore()
		h := handler.NewHandler(memStore)
		h.RegisterRoutes(r)
		fmt.Printf("Product Catalog API (in-memory) listening on :%s\n", port)
	}

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
