package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isCountryAllowedMiddleware(t *testing.T) {
	testCases := map[string]struct {
		countryCode   string
		nextHandler   http.HandlerFunc
		expStatusCode int
	}{
		"not_allowed_country_code": {
			countryCode:   "FR",
			expStatusCode: http.StatusForbidden,
		},
		"allowed_country_code": {
			countryCode: usCountryCode,
			nextHandler: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusAccepted)
			},
			expStatusCode: http.StatusAccepted,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			request := http.Request{
				Method: http.MethodGet,
				Header: http.Header{xCountryCodeHeaderKey: []string{tc.countryCode}},
			}

			rw := httptest.NewRecorder()
			isCountryAllowedMiddleware(tc.nextHandler)(rw, &request)
			assert.Equal(t, tc.expStatusCode, rw.Code)
		})
	}
}
