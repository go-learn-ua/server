package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /cards", isCountryAllowedMiddleware(listCards(storageListCards)))
	http.HandleFunc("POST /cards", isCountryAllowedMiddleware(createCard(storageSaveCard)))

	http.HandleFunc("DELETE /cards/{id}", isCountryAllowedMiddleware(deleteCard(storageDeleteCard)))
	http.HandleFunc("PUT /cards/{id}", isCountryAllowedMiddleware(updateCard(storageUpdateCard)))

	server := &http.Server{
		Addr: ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
