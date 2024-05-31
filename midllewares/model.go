package main

type creditCard struct {
	ID             int    `json:"id"`
	Number         string `json:"number"`
	ExpirationDate string `json:"expiration_date"`
	CvvCode        int    `json:"cvv"`
	Holder         string `json:"holder"`
}
