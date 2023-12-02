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
		cardsStorage  []creditCard
		countryCode   string
		expResp       string
		expStatusCode int
		queryParams   string
	}{
		"not_allowed_country_code": {
			countryCode:   "RU",
			expStatusCode: http.StatusForbidden,
		},
		"ok_one_record_in_storage": {
			countryCode: usCountryCode,
			cardsStorage: []creditCard{
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			expResp:       `[{"id":2983,"number":"4263982640269299","expiration_date":"21 січня 2023р","cvv":123,"holder":"Іванко"}]`,
			expStatusCode: http.StatusOK,
		},
		"equal filter by holder": {
			countryCode: usCountryCode,
			queryParams: "holder=Іванко",
			cardsStorage: []creditCard{
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
			},
			expResp:       `[{"id":2983,"number":"4263982640269299","expiration_date":"21 січня 2023р","cvv":123,"holder":"Іванко"}]`,
			expStatusCode: http.StatusOK,
		},
		"case insensitive filter by holder": {
			countryCode: usCountryCode,
			queryParams: "holder=івАнко",
			cardsStorage: []creditCard{
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Петрик",
				},
			},
			expResp:       `[{"id":2983,"number":"4263982640269299","expiration_date":"21 січня 2023р","cvv":123,"holder":"Іванко"}]`,
			expStatusCode: http.StatusOK,
		},
		"contains filter by holder": {
			countryCode: usCountryCode,
			queryParams: "holder=чорно",
			cardsStorage: []creditCard{
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко Чорногузко",
				},
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Петрик Чорновуско",
				},
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Не Я",
				},
			},
			expResp:       `[{"id":2983,"number":"4263982640269299","expiration_date":"21 січня 2023р","cvv":123,"holder":"Іванко Чорногузко"},{"id":2983,"number":"4263982640269299","expiration_date":"21 січня 2023р","cvv":123,"holder":"Петрик Чорновуско"}]`,
			expStatusCode: http.StatusOK,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cardsStorage = tc.cardsStorage

			request := http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{RawQuery: tc.queryParams},
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
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
		countryCode       string
		expCardsStorage   []creditCard
		expBody           string
		expStatusCode     int
	}{
		"not_allowed_country_code": {
			countryCode:   "RU",
			expStatusCode: http.StatusForbidden,
		},
		"empty_body": {
			countryCode:       ukCountryCode,
			expStatusCode:     http.StatusBadRequest,
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(errMock{}),
			expCardsStorage:   nil,
			expBody:           "",
		},
		"invalid_json": {
			countryCode:       usCountryCode,
			expStatusCode:     http.StatusBadRequest,
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(strings.NewReader("")),
			expCardsStorage:   nil,
		},
		"invalid_card_number": {
			countryCode:       ukCountryCode,
			expStatusCode:     http.StatusBadRequest,
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(strings.NewReader(`{"number":"42","expiration_date":"12/43","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage:   nil,
			expBody:           "number: must be a valid credit card number.",
		},
		"invalid_card_expiration_date": {
			countryCode:       uaCountryCode,
			expStatusCode:     http.StatusBadRequest,
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"122/43","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage:   nil,
			expBody:           "expiration_date: дата не коректна.",
		},
		"invalid_card_svv": {
			countryCode:       uaCountryCode,
			expStatusCode:     http.StatusBadRequest,
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":12223,"holder":"Іванко"}`)),
			expCardsStorage:   nil,
			expBody:           "cvv: must be no greater than 999.",
		},
		"invalid_card_holder": {
			countryCode:       uaCountryCode,
			expStatusCode:     http.StatusBadRequest,
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":123,"holder":"І"}`)),
			expCardsStorage:   nil,
			expBody:           "holder: the length must be between 5 and 50.",
		},
		"ok_empty_storage": {
			countryCode:       uaCountryCode,
			expStatusCode:     http.StatusCreated,
			setupCardsStorage: nil,
			requestBody:       io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
		"ok_with_records_in_storage": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusCreated,
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        333,
					Holder:         "Світланка",
				},
			},
			requestBody: io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":123,"holder":"Іванко"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
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
				Method: http.MethodPost,
				Body:   tc.requestBody,
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
			}

			rw := httptest.NewRecorder()
			cards(rw, &request)

			body := rw.Body.String()
			assert.Equal(t, tc.expBody, body)
			assert.Equal(t, tc.expStatusCode, rw.Code)
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

func Test_CardPut(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		path              string
		countryCode       string
		requestBody       io.ReadCloser
		expCardsStorage   []creditCard
		expBody           string
		expStatusCode     int
	}{
		"not_allowed_country_code": {
			countryCode:   "RU",
			expStatusCode: http.StatusForbidden,
		},
		"id_is_not_provided": {
			setupCardsStorage: nil,
			path:              "/cards/",
			countryCode:       ukCountryCode,
			requestBody:       nil,
			expCardsStorage:   nil,
			expStatusCode:     http.StatusNotFound,
		},
		"incorrect_id_type": {
			setupCardsStorage: nil,
			path:              "/cards/yura",
			countryCode:       ukCountryCode,
			requestBody:       nil,
			expCardsStorage:   nil,
			expStatusCode:     http.StatusBadRequest,
		},
		"validation_errors": {
			path:          "/cards/2",
			countryCode:   ukCountryCode,
			requestBody:   io.NopCloser(strings.NewReader(`{"number":"9","expiration_date":"завтра","cvv":3,"holder":"А"}`)),
			expStatusCode: http.StatusBadRequest,
			expBody:       "cvv: must be no less than 100; expiration_date: дата не коректна; holder: the length must be between 5 and 50; number: must be a valid credit card number.",
		},
		"empty_body": {
			setupCardsStorage: nil,
			countryCode:       ukCountryCode,
			path:              "/cards/1",
			requestBody:       io.NopCloser(errMock{}),
			expStatusCode:     http.StatusBadRequest,
			expCardsStorage:   nil,
		},
		"invalid_json": {
			setupCardsStorage: nil,
			countryCode:       ukCountryCode,
			path:              "/cards/1",
			requestBody:       io.NopCloser(strings.NewReader("")),
			expStatusCode:     http.StatusBadRequest,
			expCardsStorage:   nil,
		},
		"success": {
			countryCode:   ukCountryCode,
			expStatusCode: http.StatusOK,
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			path:        "/cards/2",
			requestBody: io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":337,"holder":"Петро"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        337,
					Holder:         "Петро",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
		"record_not_found": {
			countryCode:   ukCountryCode,
			expStatusCode: http.StatusNotFound,
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			path:        "/cards/5",
			requestBody: io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":337,"holder":"Петро"}`)),
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        223,
					Holder:         "Петрик",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
					CvvCode:        333,
					Holder:         "Світланка",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
					ExpirationDate: "12/43",
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
				Method: http.MethodPut,
				Body:   tc.requestBody,
				URL:    &url.URL{Path: tc.path},
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
			}

			rw := httptest.NewRecorder()
			card(rw, &request)

			body := rw.Body.String()
			assert.Equal(t, tc.expBody, body)
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

func Test_CardDelete(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		path              string
		countryCode       string
		expCardsStorage   []creditCard
		expStatusCode     int
	}{
		"not_allowed_country_code": {
			countryCode:   "RU",
			expStatusCode: http.StatusForbidden,
		},
		"id_is_not_provided": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusNotFound,
			path:          "/cards/",
		},
		"invalid_path_param": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusNotFound,
			path:          "/cards/oleh",
		},
		"record_not_found": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusNoContent,
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "нині",
					CvvCode:        1,
					Holder:         "Юра",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "завтра",
					CvvCode:        2,
					Holder:         "Олег",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
					ExpirationDate: "післязавтра",
					CvvCode:        3,
					Holder:         "Григорій",
				},
			},
			path: "/cards/183",
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "нині",
					CvvCode:        1,
					Holder:         "Юра",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "завтра",
					CvvCode:        2,
					Holder:         "Олег",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
					ExpirationDate: "післязавтра",
					CvvCode:        3,
					Holder:         "Григорій",
				},
			},
		},
		"success": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusNoContent,
			setupCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "нині",
					CvvCode:        1,
					Holder:         "Юра",
				},
				{
					ID:             2,
					Number:         "4263982640269299",
					ExpirationDate: "завтра",
					CvvCode:        2,
					Holder:         "Олег",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
					ExpirationDate: "післязавтра",
					CvvCode:        3,
					Holder:         "Григорій",
				},
			},
			path: "/cards/2",
			expCardsStorage: []creditCard{
				{
					ID:             1,
					Number:         "4263982640269299",
					ExpirationDate: "нині",
					CvvCode:        1,
					Holder:         "Юра",
				},
				{
					ID:             3,
					Number:         "4263982640269299",
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
				Method: http.MethodDelete,
				URL:    &url.URL{Path: tc.path},
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
			}

			rw := httptest.NewRecorder()
			card(rw, &request)

			assert.Empty(t, rw.Body.String())
			assert.Equal(t, tc.expStatusCode, rw.Code)
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

type errMock struct {
}

func (e errMock) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}
