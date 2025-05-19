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

type Request struct {
	Summary AccountSummary `json:"summary"`
	Email   string         `json:"email"`
}

type Response struct {
	Sent bool `json:"sent"`
}

func handler(ctx context.Context, input Request) (Response, error) {
	body := AccountSummary{
		TotalBalance:        input.Summary.TotalBalance,
		AverageCreditAmount: input.Summary.AverageCreditAmount,
		AverageDebitAmount:  input.Summary.AverageDebitAmount,
		Transactions:        input.Summary.Transactions,
	}
	parsedBody, parseError := json.Marshal(body)
	if parseError != nil {
		return Response{Sent: false}, parseError
	}
	senderSession := session.Must(session.NewSession())
	svc := ses.New(senderSession)
	_, _ = svc.SendEmail(&ses.SendEmailInput{
		Source: aws.String("axelsantacruzr@gmail.com"),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(input.Email)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{Data: aws.String("Job Status")},
			Body: &ses.Body{
				Text: &ses.Content{Data: aws.String(string(parsedBody))},
			},
		},
	})
	return Response{Sent: true}, nil
}

func main() {
	lambda.Start(handler)
}
