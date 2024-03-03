package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /cards", listCards)
	http.HandleFunc("POST /cards", createCard)

	http.HandleFunc("DELETE /cards/{id}", deleteCard)
	http.HandleFunc("PUT /cards/{id}", updateCard)

	server := &http.Server{
		Addr: ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
