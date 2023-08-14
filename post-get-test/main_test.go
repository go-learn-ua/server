package main

import (
	"net/http"
	"net/http/httptest"
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
		"empty_response": {
			cardsStorage: nil,
			expResp:      "null",
		},
		"ok_two_records_in_response": {
			cardsStorage: []creditCard{
				{
					Number:         1111111111111,
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					Number:         4444,
					ExpirationDate: "24 році",
					CvvCode:        1,
					Holder:         "Петрик",
				},
			},
			expResp: `[{"number":1111111111111,"expiration_date":"21 січня 2023р","cvv":123,"holder":"Іванко"},{"number":4444,"expiration_date":"24 році","cvv":1,"holder":"Петрик"}]`,
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
