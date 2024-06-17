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

type storageSaveCardFunc = func(card creditCard) error

func createCard(storageSaveCard storageSaveCardFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		err = storageSaveCard(reqCard)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

type storageListCardsFunc = func(holder string) ([]creditCard, error)

func listCards(storageListCards storageListCardsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		holder := r.URL.Query().Get("holder")

		creditCards, err := storageListCards(holder)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
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
}

type storageUpdateCardFunc = func(card creditCard) error

func updateCard(storageUpdateCard storageUpdateCardFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

type storageDeleteCardFunc func(id int) error

func deleteCard(storageDeleteCard storageDeleteCardFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = storageDeleteCard(id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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
