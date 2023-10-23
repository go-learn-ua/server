package main

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var cardsStorage []creditCard

type creditCard struct {
	ID             int    `json:"id"`
	Number         string `json:"number"`
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
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	switch r.Method {
	case "DELETE":
		for i := range cardsStorage {
			if cardsStorage[i].ID == id {
				cardsStorage = append(cardsStorage[:i], cardsStorage[i+1:]...)
				break
			}
		}
	case "PUT":
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

		err = validate(reqCard)
		if err != nil {
			w.Write([]byte(err.Error()))
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

		err = validate(reqCard)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		lastID := 0
		for _, card := range cardsStorage {
			if card.ID > lastID {
				lastID = card.ID
			}
		}

		reqCard.ID = lastID + 1
		cardsStorage = append(cardsStorage, reqCard)
	default:
		w.Write([]byte("Метод не підтримується!"))
	}
}

func validate(card creditCard) error {
	return validation.ValidateStruct(&card,
		validation.Field(&card.Holder, validation.Required, validation.Length(5, 50)),
		validation.Field(&card.CvvCode, validation.Required, validation.Min(100), validation.Max(999)),
		validation.Field(&card.Number, validation.Required, is.CreditCard),
		validation.Field(&card.ExpirationDate, validation.Required, validation.
			Match(regexp.MustCompile("^(0[1-9]|1[0-2])\\/[0-9]{2}$")).
			Error("дата не коректна")),
	)
}
