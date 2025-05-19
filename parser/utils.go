package main

import (
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"strconv"
	"strings"
)

type Transaction struct {
	Move          string  `json:"move_type"`
	TransactionId int     `json:"transaction_id"`
	Month         string  `json:"month"`
	Day           int     `json:"day"`
	Amount        float64 `json:"amount"`
}

func getTransactionType(value string) string {
	if value[0] == '-' {
		return "debit"
	} else {
		return "credit"
	}
}

func getTransactionAmount(value string) float64 {
	parsedValue, err := strconv.ParseFloat(value[1:], 64)
	if err != nil {
		return 0.0
	}
	return parsedValue
}

func getTransactionMonth(value string) string {
	var months = getMonths()
	month, _ := strconv.Atoi(strings.Split(value, "/")[0])
	return months[month]
}

func getTransactionDay(value string) int {
	parsedDay, _ := strconv.Atoi(strings.Split(value, "/")[1])
	return parsedDay
}

func getMonths() map[int]string {
	return map[int]string{
		1:  "January",
		2:  "February",
		3:  "March",
		4:  "April",
		5:  "May",
		6:  "June",
		7:  "July",
		8:  "August",
		9:  "September",
		10: "October",
		11: "November",
		12: "December",
	}
}

func GetFileFromS3(filename string) (io.ReadCloser, error) {
	client := s3.New(session.Must(session.NewSession()))
	response, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("stori"),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func ReadCsv(filename string, fetchFn func(string) (io.ReadCloser, error)) ([]Transaction, error) {
	s3File, err := fetchFn(filename)
	if err != nil {
		return nil, err
	}
	defer s3File.Close()
	reader := csv.NewReader(s3File)
	records, _ := reader.ReadAll()
	if len(records) == 0 {
		return nil, fmt.Errorf("csv is empty, unable to parse")
	}
	var transactions []Transaction
	for i, record := range records {
		if i == 0 {
			continue
		}
		transactionId, _ := strconv.Atoi(record[0])
		transactions = append(transactions, Transaction{
			Move:          getTransactionType(record[2]),
			Amount:        getTransactionAmount(record[2]),
			Month:         getTransactionMonth(record[1]),
			Day:           getTransactionDay(record[1]),
			TransactionId: transactionId,
		})
	}
	return transactions, nil
}
