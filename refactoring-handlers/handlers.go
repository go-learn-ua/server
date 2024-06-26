package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var cardsStorage []creditCard

type creditCard struct {
	ID             int    `json:"id"`
	Number         string `json:"number"`
	ExpirationDate string `json:"expiration_date"`
	CvvCode        int    `json:"cvv"`
	Holder         string `json:"holder"`
}

func listCards(w http.ResponseWriter, r *http.Request) {
	if isCountryAllowed(r.Header) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	creditCards := make([]creditCard, 0)
	holder := r.URL.Query().Get("holder")
	if holder != "" {
		for _, card := range cardsStorage {
			if strings.Contains(strings.ToLower(card.Holder), strings.ToLower(holder)) {
				creditCards = append(creditCards, card)
			}
		}
	} else {
		creditCards = cardsStorage
	}

	resp, err := json.Marshal(creditCards)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func createCard(w http.ResponseWriter, r *http.Request) {
	if isCountryAllowed(r.Header) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var reqCard creditCard
	err = json.Unmarshal(body, &reqCard)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate(reqCard)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
	w.WriteHeader(http.StatusCreated)
}

func updateCard(w http.ResponseWriter, r *http.Request) {
	if isCountryAllowed(r.Header) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var reqCard creditCard
	err = json.Unmarshal(body, &reqCard)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate(reqCard)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqCard.ID = id

	for i := range cardsStorage {
		if cardsStorage[i].ID == id {
			cardsStorage[i] = reqCard
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deleteCard(w http.ResponseWriter, r *http.Request) {
	if isCountryAllowed(r.Header) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for i := range cardsStorage {
		if cardsStorage[i].ID == id {
			cardsStorage = append(cardsStorage[:i], cardsStorage[i+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
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

const xCountryCodeHeaderKey = "X-Country-Code"
const uaCountryCode = "UA"
const usCountryCode = "US"
const ukCountryCode = "UK"

func isCountryAllowed(header http.Header) bool {
	code := header.Get(xCountryCodeHeaderKey)

	switch code {
	case uaCountryCode, usCountryCode, ukCountryCode:
		return true
	}
	return false
}
