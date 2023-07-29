package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var cardsStorage []creditCard

type creditCard struct {
	Number         int    `json:"number"`
	ExpirationDate string `json:"expiration_date"`
	CvvCode        int    `json:"cvv"`
	Holder         string `json:"holder"`
}

func main() {
	http.HandleFunc("/cards", cards)
	server := &http.Server{
		Addr: ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func cards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		resp, err := json.Marshal(cardsStorage)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Write(resp)
	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		r.Body.Close()

		var reqCard creditCard
		err = json.Unmarshal(body, &reqCard)
		if err != nil {
			fmt.Println(err)
			return
		}

		cardsStorage = append(cardsStorage, reqCard)
	default:
		w.Write([]byte("Метод не підтримується!"))
	}
}
