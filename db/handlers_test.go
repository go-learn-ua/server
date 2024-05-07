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
		setupStorageMock storageListCardsFunc

		countryCode   string
		expResp       string
		expStatusCode int
		queryParams   string
	}{
		"success": {
			countryCode: usCountryCode,
			setupStorageMock: func(holder string) ([]creditCard, error) {
				return []creditCard{
					{
						ID:             2983,
						Number:         "4263982640269299",
						ExpirationDate: "21 січня 2023р",
						CvvCode:        123,
						Holder:         "Іванко",
					},
				}, nil
			},
			expResp:       `[{"id":2983,"number":"4263982640269299","expiration_date":"21 січня 2023р","cvv":123,"holder":"Іванко"}]`,
			expStatusCode: http.StatusOK,
		},
		"internal_server_error": {
			countryCode: usCountryCode,
			setupStorageMock: func(holder string) ([]creditCard, error) {
				return nil, assert.AnError
			},
			expStatusCode: http.StatusInternalServerError,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			request := http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{RawQuery: tc.queryParams},
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
			}

			rw := httptest.NewRecorder()
			listCards(tc.setupStorageMock)(rw, &request)

			resp := rw.Body.String()
			expResp := tc.expResp
			assert.Equal(t, expResp, resp)
		})
	}
}

func Test_CardsPost(t *testing.T) {
	testCases := map[string]struct {
		setupStorageMock storageSaveCardFunc
		requestBody      io.ReadCloser
		countryCode      string
		expBody          string
		expStatusCode    int
	}{
		"empty_body": {
			countryCode:   ukCountryCode,
			expStatusCode: http.StatusBadRequest,
			requestBody:   io.NopCloser(errMock{}),
			expBody:       "",
		},
		"invalid_json": {
			countryCode:   usCountryCode,
			expStatusCode: http.StatusBadRequest,
			requestBody:   io.NopCloser(strings.NewReader("")),
		},
		"invalid_card_number": {
			countryCode:   ukCountryCode,
			expStatusCode: http.StatusBadRequest,
			requestBody:   io.NopCloser(strings.NewReader(`{"number":"42","expiration_date":"12/43","cvv":123,"holder":"Іванко"}`)),
			expBody:       "number: must be a valid credit card number.",
		},
		"invalid_card_expiration_date": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusBadRequest,
			requestBody:   io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"122/43","cvv":123,"holder":"Іванко"}`)),
			expBody:       "expiration_date: дата не коректна.",
		},
		"invalid_card_svv": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusBadRequest,
			requestBody:   io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":12223,"holder":"Іванко"}`)),
			expBody:       "cvv: must be no greater than 999.",
		},
		"invalid_card_holder": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusBadRequest,
			requestBody:   io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":123,"holder":"І"}`)),
			expBody:       "holder: the length must be between 5 and 50.",
		},
		"success": {
			countryCode:      uaCountryCode,
			expStatusCode:    http.StatusCreated,
			requestBody:      io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":123,"holder":"Іванко"}`)),
			setupStorageMock: func(card creditCard) error { return nil },
		},
		"internal_server_error": {
			countryCode:      uaCountryCode,
			expStatusCode:    http.StatusInternalServerError,
			requestBody:      io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":123,"holder":"Іванко"}`)),
			setupStorageMock: func(card creditCard) error { return assert.AnError },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			request := http.Request{
				Method: http.MethodPost,
				Body:   tc.requestBody,
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
			}

			rw := httptest.NewRecorder()
			createCard(tc.setupStorageMock)(rw, &request)

			body := rw.Body.String()
			assert.Equal(t, tc.expBody, body)
			assert.Equal(t, tc.expStatusCode, rw.Code)
		})
	}
}

func Test_CardPut(t *testing.T) {
	testCases := map[string]struct {
		setupStorageMock storageUpdateCardFunc
		cardID           string
		countryCode      string
		requestBody      io.ReadCloser

		expBody       string
		expStatusCode int
	}{
		"incorrect_id_type": {
			cardID:        "yura",
			countryCode:   ukCountryCode,
			requestBody:   nil,
			expStatusCode: http.StatusBadRequest,
		},
		"validation_errors": {
			cardID:        "2",
			countryCode:   ukCountryCode,
			requestBody:   io.NopCloser(strings.NewReader(`{"number":"9","expiration_date":"завтра","cvv":3,"holder":"А"}`)),
			expStatusCode: http.StatusBadRequest,
			expBody:       "cvv: must be no less than 100; expiration_date: дата не коректна; holder: the length must be between 5 and 50; number: must be a valid credit card number.",
		},
		"empty_body": {
			countryCode:   ukCountryCode,
			cardID:        "1",
			requestBody:   io.NopCloser(errMock{}),
			expStatusCode: http.StatusBadRequest,
		},
		"invalid_json": {
			countryCode:   ukCountryCode,
			cardID:        "1",
			requestBody:   io.NopCloser(strings.NewReader("")),
			expStatusCode: http.StatusBadRequest,
		},
		"success": {
			countryCode:   ukCountryCode,
			expStatusCode: http.StatusOK,
			setupStorageMock: func(card creditCard) error {
				return nil
			},
			cardID:      "2",
			requestBody: io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":337,"holder":"Петро"}`)),
		},
		"record_not_found": {
			countryCode:   ukCountryCode,
			expStatusCode: http.StatusNotFound,
			setupStorageMock: func(card creditCard) error {
				return errCreditCardNotFound
			},
			cardID:      "5",
			requestBody: io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":337,"holder":"Петро"}`)),
		},
		"internal_server_error": {
			countryCode:   ukCountryCode,
			expStatusCode: http.StatusInternalServerError,
			setupStorageMock: func(card creditCard) error {
				return assert.AnError
			},
			cardID:      "5",
			requestBody: io.NopCloser(strings.NewReader(`{"number":"4263982640269299","expiration_date":"12/43","cvv":337,"holder":"Петро"}`)),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			request := http.Request{
				Method: http.MethodPut,
				Body:   tc.requestBody,
				URL:    &url.URL{Path: "/cards/{id}"},
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
			}

			request.SetPathValue("id", tc.cardID)

			rw := httptest.NewRecorder()
			updateCard(tc.setupStorageMock)(rw, &request)

			body := rw.Body.String()
			assert.Equal(t, tc.expBody, body)
			assert.Equal(t, tc.expStatusCode, rw.Code)
		})
	}
}

func Test_CardDelete(t *testing.T) {
	testCases := map[string]struct {
		setupStorageMock storageDeleteCardFunc
		cardID           string
		countryCode      string
		expStatusCode    int
	}{
		"invalid_path_param": {
			countryCode:   uaCountryCode,
			expStatusCode: http.StatusNotFound,
			cardID:        "oleh",
		},
		"success": {
			countryCode:      uaCountryCode,
			expStatusCode:    http.StatusNoContent,
			cardID:           "12",
			setupStorageMock: func(id int) error { return nil },
		},
		"internal_server_error": {
			countryCode:      uaCountryCode,
			expStatusCode:    http.StatusInternalServerError,
			cardID:           "5",
			setupStorageMock: func(id int) error { return assert.AnError },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			request := http.Request{
				Method: http.MethodDelete,
				URL:    &url.URL{Path: "/cards/{id}"},
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
			}
			request.SetPathValue("id", tc.cardID)

			rw := httptest.NewRecorder()
			deleteCard(tc.setupStorageMock)(rw, &request)

			assert.Empty(t, rw.Body.String())
			assert.Equal(t, tc.expStatusCode, rw.Code)
		})
	}
}

type errMock struct {
}

func (e errMock) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}
