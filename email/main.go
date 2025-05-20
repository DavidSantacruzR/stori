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
	Sender  string         `json:"sender"`
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
	parsedBody, parsingErr := json.Marshal(body)
	if parsingErr != nil {
		return Response{false}, parsingErr
	}
	emailConfig := &ses.SendTemplatedEmailInput{
		Source: aws.String(input.Sender),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(input.Email)},
		},
		Template:     aws.String("StoriSummaryTemplate"),
		TemplateData: aws.String(string(parsedBody)),
	}
	senderSession := session.Must(session.NewSession())
	svc := ses.New(senderSession)
	_, _ = svc.SendTemplatedEmail(emailConfig)
	return Response{Sent: true}, nil
}

func main() {
	lambda.Start(handler)
}
