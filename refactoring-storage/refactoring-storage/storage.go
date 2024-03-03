package main

import (
	"errors"
	"strings"
)

var errCreditCardNotFound = errors.New("credit card not found")

var cardsStorage []creditCard

func storageSaveCard(card creditCard) {
	lastID := 0
	for _, card := range cardsStorage {
		if card.ID > lastID {
			lastID = card.ID
		}
	}

	card.ID = lastID + 1
	cardsStorage = append(cardsStorage, card)
}

func storageListCards(holder string) []creditCard {
	creditCards := make([]creditCard, 0)
	if holder != "" {
		for _, card := range cardsStorage {
			if strings.Contains(strings.ToLower(card.Holder), strings.ToLower(holder)) {
				creditCards = append(creditCards, card)
			}
		}
	} else {
		creditCards = cardsStorage
	}

	return creditCards
}

func storageUpdateCard(card creditCard) error {
	for i := range cardsStorage {
		if cardsStorage[i].ID == card.ID {
			cardsStorage[i] = card
			return nil
		}
	}

	return errCreditCardNotFound
}

func storageDeleteCard(id int) {
	for i := range cardsStorage {
		if cardsStorage[i].ID != id {
			continue
		}

		cardsStorage = append(cardsStorage[:i], cardsStorage[i+1:]...)
		return
	}
}
