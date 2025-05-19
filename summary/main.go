package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Transactions []Transaction `json:"transactions"`
	Email        string        `json:"email"`
}

type Response struct {
	Summary AccountSummary `json:"summary"`
	Email   string         `json:"email"`
}

func handler(ctx context.Context, input Request) (Response, error) {
	summary, err := GetAccountSummary(input.Transactions)
	if err != nil {
		return Response{}, err
	}
	return Response{Summary: summary, Email: input.Email}, nil
}

func main() {
	lambda.Start(handler)
}
