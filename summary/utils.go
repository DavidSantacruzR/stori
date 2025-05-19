package main

type AccountSummary struct {
	TotalBalance        float64        `json:"total_balance"`
	AverageCreditAmount float64        `json:"average_credit_amount"`
	AverageDebitAmount  float64        `json:"average_debit_amount"`
	Transactions        []MonthSummary `json:"monthly_summary"`
}

type MonthSummary struct {
	Month                string `json:"month"`
	NumberOfTransactions int    `json:"number_of_transactions"`
}

type Transaction struct {
	Move          string  `json:"move_type"`
	TransactionId int     `json:"transaction_id"`
	Month         string  `json:"month"`
	Day           int     `json:"day"`
	Amount        float64 `json:"amount"`
}

func GetAccountSummary(transactions []Transaction) (AccountSummary, error) {
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
	}, nil
}
