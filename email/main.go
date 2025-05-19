package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type MonthSummary struct {
	Month                string `json:"month"`
	NumberOfTransactions int    `json:"number_of_transactions"`
}

type AccountSummary struct {
	TotalBalance        float64        `json:"total_balance"`
	AverageCreditAmount float64        `json:"average_credit_amount"`
	AverageDebitAmount  float64        `json:"average_debit_amount"`
	Transactions        []MonthSummary `json:"monthly_summary"`
}

func handler(ctx context.Context, input string) error {
	var summary = AccountSummary{}
	err := json.Unmarshal([]byte(input), &summary)
	if err != nil {
		return err
	}
	body := AccountSummary{
		TotalBalance:        summary.TotalBalance,
		AverageCreditAmount: summary.AverageCreditAmount,
		AverageDebitAmount:  summary.AverageDebitAmount,
		Transactions:        summary.Transactions,
	}
	parsedBody, parseError := json.Marshal(body)
	if parseError != nil {
		return parseError
	}
	senderSession := session.Must(session.NewSession())
	svc := ses.New(senderSession)
	_, err = svc.SendEmail(&ses.SendEmailInput{
		Source: aws.String("axelsantacruzr@gmail.com"),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(ctx.Value("email").(string))},
		},
		Message: &ses.Message{
			Subject: &ses.Content{Data: aws.String("Job Status")},
			Body: &ses.Body{
				Text: &ses.Content{Data: aws.String(string(parsedBody))},
			},
		},
	})
	return err
}

func main() {
	lambda.Start(handler)
}
