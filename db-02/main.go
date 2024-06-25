package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	connStr := "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"

	var err error
	cardsStorage, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
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
