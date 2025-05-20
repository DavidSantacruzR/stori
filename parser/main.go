package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Email    string `json:"email"`
	Sender   string `json:"sender"`
	Filename string `json:"filename"`
}

type Response struct {
	Transactions []Transaction `json:"transactions"`
	Email        string        `json:"email"`
	Sender       string        `json:"sender"`
	Filename     string        `json:"filename"`
}

func handler(ctx context.Context, input Request) (Response, error) {
	parsedCsv, csvErr := ReadCsv(input.Filename, GetFileFromS3)
	if csvErr != nil {
		return Response{}, csvErr
	}
	ctx = context.WithValue(ctx, "email", input.Email)
	return Response{Transactions: parsedCsv, Email: input.Email, Sender: input.Sender, Filename: input.Filename}, nil
}

func main() {
	lambda.Start(handler)
}
