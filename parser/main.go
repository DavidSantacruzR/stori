package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Email    string `json:"email"`
	Filename string `json:"filename"`
}

func handler(ctx context.Context, input Request) (string, error) {
	parsedCsv, csvErr := ReadCsv(input.Filename, GetFileFromS3)
	if csvErr != nil {
		return "", csvErr
	}
	jsonResult, jsonParsingErr := json.Marshal(parsedCsv)
	if jsonParsingErr != nil {
		return "", jsonParsingErr
	}
	ctx = context.WithValue(ctx, "email", input.Email)
	return string(jsonResult), nil
}

func main() {
	lambda.Start(handler)
}
