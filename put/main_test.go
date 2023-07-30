package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
					ID:             1,
					Number:         3911391723597698,
					ExpirationDate: "20 липня 2031р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			expResp: `[{"id":1,"number":3911391723597698,"expiration_date":"20 липня 2031р","cvv":123,"holder":"Іванко"}]`,
		},
		"empty_response": {
			expResp: "null",
		},
		"ok_two_records_in_response": {
			cardsStorage: []creditCard{
				{
					ID:             1,
					Number:         3911391723597698,
					ExpirationDate: "20 липня 2031р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					ID:             2,
					Number:         4444,
					ExpirationDate: "20 липня 2024р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
			},
			expResp: `[{"id":1,"number":3911391723597698,"expiration_date":"20 липня 2031р","cvv":123,"holder":"Іванко"},{"id":2,"number":4444,"expiration_date":"20 липня 2024р","cvv":123,"holder":"Петрик"}]`,
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
			assert.Equal(t, tc.expResp, resp)
		})
	}

	cardsStorage = nil
}

func Test_CardsPost(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		requestBody       io.ReadCloser
		expCardsStorage   []creditCard
	}{
		"empty_body": {
			requestBody:     io.NopCloser(readerErrMock{}),
			expCardsStorage: nil,
		},
		"invalid_json": {
			requestBody:     io.NopCloser(strings.NewReader(``)),
			expCardsStorage: nil,
		},
		"ok_empty_storage": {
			requestBody: io.NopCloser(strings.NewReader(`{"number":3911391723597698,"expiration_date":"20 липня 2031р","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         3911391723597698,
					ExpirationDate: "20 липня 2031р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
		"ok_with_records_in_storage": {
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         1111,
					ExpirationDate: "20 липня 2032р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         2222,
					ExpirationDate: "20 липня 2033р",
					CvvCode:        333,
					Holder:         "Світланка",
				},
			},
			requestBody: io.NopCloser(strings.NewReader(`{"number":3911391723597698,"expiration_date":"20 липня 2031р","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             3,
					Number:         3911391723597698,
					ExpirationDate: "20 липня 2031р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					ID:             1,
					Number:         1111,
					ExpirationDate: "20 липня 2032р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         2222,
					ExpirationDate: "20 липня 2033р",
					CvvCode:        333,
					Holder:         "Світланка",
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

			response := rw.Body.String()
			assert.Empty(t, response)
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}

	cardsStorage = nil
}

func Test_CardPut(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		path              string
		requestBody       io.ReadCloser
		expCardsStorage   []creditCard
	}{
		"id_is_not_provided": {
			path: "/cards/",
		},
		"incorrect_id_type": {
			path: "/cards/yura",
		},
		"empty_body": {
			requestBody:     io.NopCloser(readerErrMock{}),
			path:            "/cards/1",
			expCardsStorage: nil,
		},
		"invalid_json": {
			requestBody:     io.NopCloser(strings.NewReader(``)),
			path:            "/cards/1",
			expCardsStorage: nil,
		},
		"success": {
			setupCardsStorage: []creditCard{
				{
					ID:             3,
					Number:         3911391723597698,
					ExpirationDate: "20 липня 2031р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					ID:             1,
					Number:         1111,
					ExpirationDate: "20 липня 2032р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         2222,
					ExpirationDate: "20 липня 2033р",
					CvvCode:        333,
					Holder:         "Світланка",
				},
			},
			path:        "/cards/2",
			requestBody: io.NopCloser(strings.NewReader(`{"number":17,"expiration_date":"завтра","cvv":7,"holder":"Петро"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             3,
					Number:         3911391723597698,
					ExpirationDate: "20 липня 2031р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					ID:             1,
					Number:         1111,
					ExpirationDate: "20 липня 2032р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         17,
					ExpirationDate: "завтра",
					CvvCode:        7,
					Holder:         "Петро",
				},
			},
		},
		"record not found": {
			setupCardsStorage: []creditCard{
				{
					ID:             3,
					Number:         3911391723597698,
					ExpirationDate: "20 липня 2031р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					ID:             1,
					Number:         1111,
					ExpirationDate: "20 липня 2032р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         2222,
					ExpirationDate: "20 липня 2033р",
					CvvCode:        333,
					Holder:         "Світланка",
				},
			},
			path:        "/cards/5",
			requestBody: io.NopCloser(strings.NewReader(`{"number":17,"expiration_date":"завтра","cvv":7,"holder":"Петро"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             3,
					Number:         3911391723597698,
					ExpirationDate: "20 липня 2031р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					ID:             1,
					Number:         1111,
					ExpirationDate: "20 липня 2032р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         2222,
					ExpirationDate: "20 липня 2033р",
					CvvCode:        333,
					Holder:         "Світланка",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cardsStorage = tc.setupCardsStorage

			request := http.Request{
				Method: "PUT",
				Body:   tc.requestBody,
				URL:    &url.URL{Path: tc.path},
			}

			rw := httptest.NewRecorder()
			card(rw, &request)

			response := rw.Body.String()
			assert.Empty(t, response)
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

type readerErrMock struct{}

func (r readerErrMock) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}
