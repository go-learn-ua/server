package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type storageSaveCardFunc = func(card creditCard)

func createCard(storageSaveCard storageSaveCardFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		storageSaveCard(reqCard)
		w.WriteHeader(http.StatusCreated)
	}
}

type storageListCardsFunc = func(holder string) []creditCard

func listCards(storageListCards storageListCardsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isCountryAllowed(r.Header) == false {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		holder := r.URL.Query().Get("holder")

		creditCards := storageListCards(holder)
		resp, err := json.Marshal(creditCards)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

type storageUpdateCardFunc = func(card creditCard) error

func updateCard(storageUpdateCard storageUpdateCardFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("start")
		if isCountryAllowed(r.Header) == false {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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

		reqCard.ID = id

		err = storageUpdateCard(reqCard)
		if err == errCreditCardNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

type storageDeleteCardFunc func(id int)

func deleteCard(storageDeleteCard storageDeleteCardFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isCountryAllowed(r.Header) == false {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		storageDeleteCard(id)
		w.WriteHeader(http.StatusNoContent)
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
