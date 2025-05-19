package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func mockGetFileFromS3(filename string) (io.ReadCloser, error) {
	csvData := `
	Id,Date,Transaction
	0,7/15,+60.5
	1,7/28,-10.3
	2,8/2,-20.46
	3,8/13,+10`
	csvData = strings.TrimSpace(csvData)
	csvData = strings.ReplaceAll(csvData, "\t", "")
	return io.NopCloser(strings.NewReader(csvData)), nil
}

func mockGetFileFromS3Empty(filename string) (io.ReadCloser, error) {
	csvData := ``
	return io.NopCloser(strings.NewReader(csvData)), nil
}

func TestParseCsvFileSuccess(t *testing.T) {
	file, _ := ReadCsv("transactions.csv", mockGetFileFromS3)
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
	fmt.Println(file)
	fmt.Println(expectedResult)
	if !reflect.DeepEqual(file, expectedResult) {
		t.Error("Parsed test csv does not match the expected format.")
	}
}

func TestParseCsvFileError(t *testing.T) {
	_, err := ReadCsv("transactionsEmpty.csv", mockGetFileFromS3Empty)
	expectedError := "csv is empty, unable to parse"
	if err.Error() != expectedError {
		t.Error("Error does not match the expected error message")
	}
}
