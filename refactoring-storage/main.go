package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /cards", listCards(storageListCards))
	http.HandleFunc("POST /cards", createCard(storageSaveCard))

	http.HandleFunc("DELETE /cards/{id}", deleteCard(storageDeleteCard))
	http.HandleFunc("PUT /cards/{id}", updateCard(storageUpdateCard))

	server := &http.Server{
		Addr: ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
