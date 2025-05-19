package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, input string) (string, error) {
	summary, err := GetAccountSummary(input)
	if err != nil {
		return "", err
	}
	jsonResult, marshalErr := json.Marshal(summary)
	if marshalErr != nil {
		return "", marshalErr
	}
	return string(jsonResult), nil
}

func main() {
	lambda.Start(handler)
}
