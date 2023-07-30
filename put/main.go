package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var cardsStorage []creditCard

type creditCard struct {
	ID             int    `json:"id"`
	Number         int    `json:"number"`
	ExpirationDate string `json:"expiration_date"`
	CvvCode        int    `json:"cvv"`
	Holder         string `json:"holder"`
}

func main() {
	http.HandleFunc("/cards", cards)
	http.HandleFunc("/cards/", card)
	server := &http.Server{
		Addr: ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func card(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		path := strings.Split(r.URL.Path, "/")
		id, err := strconv.Atoi(path[2])
		if err != nil {
			fmt.Println(err)
			return
		}

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

		reqCard.ID = id
		for i := range cardsStorage {
			if cardsStorage[i].ID == id {
				cardsStorage[i] = reqCard
			}
		}
	default:
		w.Write([]byte("Метод не підтримується!"))
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

		lastID := 0
		for _, c := range cardsStorage {
			if c.ID > lastID {
				lastID = c.ID
			}
		}

		reqCard.ID = lastID + 1
		cardsStorage = append(cardsStorage, reqCard)
	default:
		w.Write([]byte("Метод не підтримується!"))
	}
}
