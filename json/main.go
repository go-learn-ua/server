package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type wallet struct {
	CreditCards []creditCard `json:"credit_cards"`
}

type creditCard struct {
	Number         int    `json:"number"`
	ExpirationDate string `json:"expiration_date"`
	CvvCode        int    `json:"-"`
	Holder         string `json:"holder,omitempty"`
}

var ivankoCreditCard = creditCard{
	Number:         4224580604397698,
	ExpirationDate: "20 липня 2031р",
	CvvCode:        230,
	Holder:         "",
}

var ivankoMonoCreditCard = creditCard{
	Number:         4114391864397698,
	ExpirationDate: "20 липня 2025р",
	CvvCode:        892,
}

var ivankoPrivatCreditCard = creditCard{
	Number:         3911391723597698,
	ExpirationDate: "8 вересня 2025р",
	CvvCode:        777,
	Holder:         "Іванко",
}

var myWallet = wallet{
	CreditCards: []creditCard{
		ivankoMonoCreditCard,
		ivankoPrivatCreditCard,
		ivankoCreditCard,
	},
}

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
	result, err := json.Marshal(myWallet)
	if err != nil {
		panic(err)
	}

	var newWallet wallet
	err = json.Unmarshal(result, &newWallet)
	if err != nil {
		panic(err)
	}

	newWallet.CreditCards[0].Holder = "Юра"
	fmt.Println(newWallet.CreditCards[0].Holder)
	fmt.Println(myWallet.CreditCards[0].Holder)

	w.Write(result)
}
