package main

import "net/http"

func isCountryAllowedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isCountryAllowed(r.Header) == false {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
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
