package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Transaction struct {
	Move          string
	TransactionId int
	Month         string
	Day           int
	Amount        float64
}

type AccountSummary struct {
	TotalBalance        float64
	AverageCreditAmount float64
	AverageDebitAmount  float64
	Transactions        []MonthSummary
}

type MonthSummary struct {
	Month                string
	NumberOfTransactions int
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

func ReadCsv(filename string) ([]Transaction, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
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

func GetAccountSummary(filename string) AccountSummary {
	transactions, err := ReadCsv(filename)
	if err != nil {
		log.Fatal(err)
	}
	totalCreditAmount := 0.0
	totalCreditTransactions := 0.0
	totalDebitAmount := 0.0
	totalDebitTransactions := 0.0
	monthlySummary := make(map[string]int)
	for _, transaction := range transactions {
		if transaction.Move == "debit" {
			totalDebitAmount = totalDebitAmount + transaction.Amount
			totalDebitTransactions++
		} else {
			totalCreditAmount = totalCreditAmount + transaction.Amount
			totalCreditTransactions++
		}
		monthlySummary[transaction.Month]++
	}
	var parsedMonthlySummary []MonthSummary
	for month, count := range monthlySummary {
		parsedMonthlySummary = append(parsedMonthlySummary, MonthSummary{
			Month:                month,
			NumberOfTransactions: count,
		})
	}
	return AccountSummary{
		totalCreditAmount - totalDebitAmount,
		totalCreditAmount / totalCreditTransactions,
		totalDebitAmount / totalDebitTransactions,
		parsedMonthlySummary,
	}
}
