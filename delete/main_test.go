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
					ID:             2983,
					Number:         1111111111111,
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			expResp: `[{"id":2983,"number":1111111111111,"expiration_date":"21 січня 2023р","cvv":123,"holder":"Іванко"}]`,
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
					ID:             1,
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
					ID:             1,
					Number:         22222,
					ExpirationDate: "24 липня 2091",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         333333,
					ExpirationDate: "8 серпня 2082",
					CvvCode:        333,
					Holder:         "Світланка",
				},
			},
			requestBody: io.NopCloser(strings.NewReader(`{"number":111111,"expiration_date":"21 липня 2025р","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         22222,
					ExpirationDate: "24 липня 2091",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         333333,
					ExpirationDate: "8 серпня 2082",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					ID:             3,
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

func Test_CardPut(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		path              string
		requestBody       io.ReadCloser
		expCardsStorage   []creditCard
	}{
		"id_is_not_provided": {
			setupCardsStorage: nil,
			path:              "/cards/",
			requestBody:       nil,
			expCardsStorage:   nil,
		},
		"incorrect_id_type": {
			setupCardsStorage: nil,
			path:              "/cards/yura",
			requestBody:       nil,
			expCardsStorage:   nil,
		},
		"empty_body": {
			setupCardsStorage: nil,
			path:              "/cards/1",
			requestBody:       io.NopCloser(errMock{}),
			expCardsStorage:   nil,
		},
		"invalid_json": {
			setupCardsStorage: nil,
			path:              "/cards/1",
			requestBody:       io.NopCloser(strings.NewReader("")),
			expCardsStorage:   nil,
		},
		"success": {
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         22222,
					ExpirationDate: "24 липня 2091",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         333333,
					ExpirationDate: "8 серпня 2082",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					ID:             3,
					Number:         111111,
					ExpirationDate: "21 липня 2025р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			path:        "/cards/2",
			requestBody: io.NopCloser(strings.NewReader(`{"number":17,"expiration_date":"завтра","cvv":7,"holder":"Петро"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         22222,
					ExpirationDate: "24 липня 2091",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         17,
					ExpirationDate: "завтра",
					CvvCode:        7,
					Holder:         "Петро",
				},
				{
					ID:             3,
					Number:         111111,
					ExpirationDate: "21 липня 2025р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
		"record_not_found": {
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         22222,
					ExpirationDate: "24 липня 2091",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         333333,
					ExpirationDate: "8 серпня 2082",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					ID:             3,
					Number:         111111,
					ExpirationDate: "21 липня 2025р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			path:        "/cards/5",
			requestBody: io.NopCloser(strings.NewReader(`{"number":17,"expiration_date":"завтра","cvv":7,"holder":"Петро"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         22222,
					ExpirationDate: "24 липня 2091",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         333333,
					ExpirationDate: "8 серпня 2082",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					ID:             3,
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
				Method: "PUT",
				Body:   tc.requestBody,
				URL:    &url.URL{Path: tc.path},
			}

			rw := httptest.NewRecorder()
			card(rw, &request)

			resp := rw.Body.String()
			assert.Empty(t, resp)
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

func Test_CardDelete(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		path              string
		expCardsStorage   []creditCard
	}{
		"id_is_not_provided": {
			path: "/cards/",
		},
		"invalid_path_param": {
			path: "/cards/oleh",
		},
		"record_not_found": {
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         11111,
					ExpirationDate: "нині",
					CvvCode:        1,
					Holder:         "Юра",
				},
				{
					ID:             2,
					Number:         222222,
					ExpirationDate: "завтра",
					CvvCode:        2,
					Holder:         "Олег",
				},
				{
					ID:             3,
					Number:         33333,
					ExpirationDate: "післязавтра",
					CvvCode:        3,
					Holder:         "Григорій",
				},
			},
			path: "/cards/183",
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         11111,
					ExpirationDate: "нині",
					CvvCode:        1,
					Holder:         "Юра",
				},
				{
					ID:             2,
					Number:         222222,
					ExpirationDate: "завтра",
					CvvCode:        2,
					Holder:         "Олег",
				},
				{
					ID:             3,
					Number:         33333,
					ExpirationDate: "післязавтра",
					CvvCode:        3,
					Holder:         "Григорій",
				},
			},
		},
		"success": {
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         11111,
					ExpirationDate: "нині",
					CvvCode:        1,
					Holder:         "Юра",
				},
				{
					ID:             2,
					Number:         222222,
					ExpirationDate: "завтра",
					CvvCode:        2,
					Holder:         "Олег",
				},
				{
					ID:             3,
					Number:         33333,
					ExpirationDate: "післязавтра",
					CvvCode:        3,
					Holder:         "Григорій",
				},
			},
			path: "/cards/2",
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         11111,
					ExpirationDate: "нині",
					CvvCode:        1,
					Holder:         "Юра",
				},
				{
					ID:             3,
					Number:         33333,
					ExpirationDate: "післязавтра",
					CvvCode:        3,
					Holder:         "Григорій",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cardsStorage = tc.setupCardsStorage

			request := http.Request{
				Method: "DELETE",
				URL:    &url.URL{Path: tc.path},
			}

			rw := httptest.NewRecorder()
			card(rw, &request)

			assert.Empty(t, rw.Body.String())
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

type errMock struct {
}

func (e errMock) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}
