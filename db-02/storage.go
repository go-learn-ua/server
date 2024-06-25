package main

import (
	"database/sql"
	"errors"
	"fmt"
)

var errCreditCardNotFound = errors.New("credit card not found")

var cardsStorage *sql.DB

func storageSaveCard(card creditCard) error {
	res, err := cardsStorage.Exec("INSERT INTO credit_cards(number, expiration_date, cvv, holder_name) VALUES ($1, $2, $3, $4)",
		card.Number, card.ExpirationDate, card.CvvCode, card.Holder)
	if err != nil {
		return fmt.Errorf("exec insert into credit cards: %w", err)
	}

	numRowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected on insert: %w", err)
	}
	if numRowsAffected != 1 {
		return fmt.Errorf("invalid number of rows affected: %d", numRowsAffected)
	}

	return nil
}

func storageListCards(holder string) ([]creditCard, error) {
	var rows *sql.Rows
	var err error

	if holder != "" {
		rows, err = cardsStorage.Query("SELECT id, number, expiration_date, cvv, holder_name FROM credit_cards WHERE LOWER(holder_name) LIKE LOWER($1)",
			"%"+holder+"%")
	} else {
		rows, err = cardsStorage.Query("SELECT id, number, expiration_date, cvv, holder_name FROM credit_cards")
	}

	if err != nil {
		return nil, fmt.Errorf("query credit cards: %w", err)
	}

	creditCards := make([]creditCard, 0)
	for rows.Next() {
		var card creditCard
		err := rows.Scan(&card.ID, &card.Number, &card.ExpirationDate, &card.CvvCode, &card.Holder)
		if err != nil {
			return nil, fmt.Errorf("scan credit card: %w", err)
		}

		creditCards = append(creditCards, card)
	}

	return creditCards, nil
}

func storageUpdateCard(card creditCard) error {
	res, err := cardsStorage.Exec("UPDATE credit_cards SET number=$1, expiration_date=$2, cvv=$3, holder_name=$4 WHERE id=$5",
		card.Number, card.ExpirationDate, card.CvvCode, card.Holder, card.ID)
	if err != nil {
		return fmt.Errorf("exec update into credit cards: %w", err)
	}

	numRowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected on update: %w", err)
	}
	if numRowsAffected != 1 {
		return errCreditCardNotFound
	}

	return nil
}

func storageDeleteCard(id int) error {
	_, err := cardsStorage.Exec("DELETE FROM credit_cards WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("exec delete from credit cards: %w", err)
	}

	return nil
}
