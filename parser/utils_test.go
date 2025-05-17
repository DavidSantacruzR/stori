package main

import (
	"reflect"
	"testing"
)

func TestParseCsvFileSuccess(t *testing.T) {
	file, _ := ReadCsv("transactions.csv")
	expectedResult := []Transaction{
		{
			TransactionId: 0,
			Month:         "July",
			Day:           15,
			Move:          "credit",
			Amount:        60.5,
		},
		{
			TransactionId: 1,
			Month:         "July",
			Day:           28,
			Move:          "debit",
			Amount:        10.3,
		},
		{
			TransactionId: 2,
			Month:         "August",
			Day:           2,
			Move:          "debit",
			Amount:        20.46,
		},
		{
			TransactionId: 3,
			Month:         "August",
			Day:           13,
			Move:          "credit",
			Amount:        10.0,
		},
	}
	if !reflect.DeepEqual(file, expectedResult) {
		t.Error("Parsed test csv does not match the expected format.")
	}
}

func TestParseCsvFileError(t *testing.T) {
	_, err := ReadCsv("transactionsEmpty.csv")
	expectedError := "csv is empty, unable to parse"
	if err.Error() != expectedError {
		t.Error("Error does not match the expected error message")
	}
}

func TestGetAccountSummary(t *testing.T) {
	value := GetAccountSummary("transactions.csv")
	monthlySummary := make([]MonthSummary, 2)
	monthlySummary = append(monthlySummary, MonthSummary{
		Month:                "July",
		NumberOfTransactions: 2,
	})
	monthlySummary = append(monthlySummary, MonthSummary{
		Month:                "August",
		NumberOfTransactions: 2,
	})
	expectedResult := AccountSummary{
		TotalBalance:        39.739999999999995,
		AverageCreditAmount: 35.25,
		AverageDebitAmount:  15.38,
		Transactions:        monthlySummary,
	}
	if reflect.DeepEqual(value, expectedResult) {
	}
}
