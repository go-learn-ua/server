package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	cardsStorage, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer cardsStorage.Close()

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(cardsStorage, "migrations"); err != nil {
		panic(err)
	}

	http.HandleFunc("GET /cards", isCountryAllowedMiddleware(listCards(storageListCards)))
	http.HandleFunc("POST /cards", isCountryAllowedMiddleware(createCard(storageSaveCard)))

	http.HandleFunc("DELETE /cards/{id}", isCountryAllowedMiddleware(deleteCard(storageDeleteCard)))
	http.HandleFunc("PUT /cards/{id}", isCountryAllowedMiddleware(updateCard(storageUpdateCard)))

	server := &http.Server{
		Addr: ":8080",
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
