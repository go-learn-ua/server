package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", api)

	server := &http.Server{
		Addr: ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func api(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	fmt.Fprintln(w, "query params:", queryParams.Encode())

	path := r.URL.Path
	path = strings.Trim(path, "/")
	pathParams := strings.Split(path, "/")
	if len(pathParams) > 0 {
		fmt.Fprintln(w, "path params:", pathParams[0])
	}

	resp, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "body:", string(resp))
}
