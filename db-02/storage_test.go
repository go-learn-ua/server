package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_SaveCard(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		card              creditCard
		expCardsStorage   []creditCard
	}{
		"ok_empty_storage": {
			setupCardsStorage: nil,
			card: creditCard{
				Number:         "4263982640269299",
				ExpirationDate: "12/43",
				CvvCode:        123,
				Holder:         "Іванко",
			},
			expCardsStorage: []creditCard{{
				ID:             1,
				Number:         "4263982640269299",
				ExpirationDate: "12/43",
				CvvCode:        123,
				Holder:         "Іванко",
			}},
		},
		"ok_with_records_in_storage": {
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
			card: creditCard{
				Number:         "4263982640269299",
				ExpirationDate: "12/43",
				CvvCode:        123,
				Holder:         "Іванко",
			},
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

			storageSaveCard(tc.card)
			assert.ElementsMatch(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

func TestStorage_ListCards(t *testing.T) {
	testCases := map[string]struct {
		cardsStorage []creditCard
		holder       string
		expResp      []creditCard
	}{
		"ok_one_record_in_storage": {
			cardsStorage: []creditCard{
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
			holder: "",
			expResp: []creditCard{
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
		"equal filter by holder": {
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
			holder: "Іванко",
			expResp: []creditCard{
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
		"case insensitive filter by holder": {
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
			holder: "івАнко",
			expResp: []creditCard{
				{
					ID:             2983,
					Number:         "4263982640269299",
					ExpirationDate: "21 січня 2023р",
					CvvCode:        123,
					Holder:         "Іванко",
				},
			},
		},
		"contains filter by holder": {
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
			holder: "чорНо",
			expResp: []creditCard{
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
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cardsStorage = tc.cardsStorage

			gotCards := storageListCards(tc.holder)
			assert.Equal(t, tc.expResp, gotCards)
		})
	}
}

func TestStorage_UpdateCard(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		card              creditCard
		expCardsStorage   []creditCard
		expErr            error
	}{
		"success": {
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
			card: creditCard{
				ID:             2,
				Number:         "4263982640269299",
				ExpirationDate: "12/43",
				CvvCode:        337,
				Holder:         "Петро",
			},
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
			expErr: nil,
		},
		"record_not_found": {
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
			card: creditCard{
				ID: 5,
			},
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
			expErr: errCreditCardNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cardsStorage = tc.setupCardsStorage

			gotErr := storageUpdateCard(tc.card)
			assert.Equal(t, tc.expErr, gotErr)
			assert.Equal(t, tc.expCardsStorage, cardsStorage)
		})
	}
}

func TestStorage_DeleteCard(t *testing.T) {
	testCases := map[string]struct {
		setupCardsStorage []creditCard
		cardID            int
		expCardsStorage   []creditCard
	}{
		"record_not_found": {
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
			cardID: 83,
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
			cardID: 2,
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

			storageDeleteCard(tc.cardID)
			assert.Equal(t, tc.expCardsStorage, cardsStorage)
		})
	}
}
