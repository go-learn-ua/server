package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CardsGet(t *testing.T) {
	testCases := map[string]struct {
		cardsStorage []creditCard
		expResp      string
	}{
		"ok_one_record_in_storage": {
			cardsStorage: []creditCard{
				{
					Number:         1111111111111,
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			expResp: `[{"number":1111111111111,"expiration_date":"21 січня 2023р","cvv":123,"holder":"Іванко"}]`,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cardsStorage = tc.cardsStorage

			request := http.Request{
				Method: "GET",
			}

			rw := httptest.NewRecorder()

			cards(rw, &request)

			resp := rw.Body.String()
			expResp := tc.expResp
			assert.Equal(t, expResp, resp)
		})
	}
}

func Test_CardsPost(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		requestBody       io.ReadCloser
		expCardsStorage   []creditCard
	}{
		"empty_body": {
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(errMock{}),
			expCardsStorage:   nil,
		},
		"invalid_json": {
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(strings.NewReader("")),
			expCardsStorage:   nil,
		},
		"ok_empty_storage": {
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(strings.NewReader(`{"number":111111,"expiration_date":"21 липня 2025р","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage: []creditCard{
				{
					Number:         111111,
					ExpirationDate: "21 липня 2025р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
		"ok_with_records_in_storage": {
			setupCardsStorage: []creditCard{
				{
					Number:         22222,
					ExpirationDate: "24 липня 2091",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					Number:         333333,
					ExpirationDate: "8 серпня 2082",
					CvvCode:        333,
					Holder:         "Світланка",
				},
			},
			requestBody: io.NopCloser(strings.NewReader(`{"number":111111,"expiration_date":"21 липня 2025р","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage: []creditCard{
				{
					Number:         22222,
					ExpirationDate: "24 липня 2091",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					Number:         333333,
					ExpirationDate: "8 серпня 2082",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					Number:         111111,
					ExpirationDate: "21 липня 2025р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cardsStorage = tc.setupCardsStorage
			request := http.Request{
				Method: "POST",
				Body:   tc.requestBody,
			}

			rw := httptest.NewRecorder()
			cards(rw, &request)

			body := rw.Body.String()
			assert.Empty(t, body)
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

type errMock struct {
}

func (e errMock) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}
