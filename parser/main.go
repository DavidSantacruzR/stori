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

func handler(ctx context.Context, input string) (string, error) {
	var request = Request{}
	requestErr := json.Unmarshal([]byte(input), &request)
	if requestErr != nil {
		return "", requestErr
	}
	parsedCsv, csvErr := ReadCsv(request.Filename)
	if csvErr != nil {
		return "", csvErr
	}
	jsonResult, jsonParsingErr := json.Marshal(parsedCsv)
	if jsonParsingErr != nil {
		return "", jsonParsingErr
	}
	ctx = context.WithValue(ctx, "email", request.Email)
	return string(jsonResult), nil
}

func main() {
	lambda.Start(handler)
}
