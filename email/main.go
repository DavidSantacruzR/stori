package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
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
	Summary  AccountSummary `json:"summary"`
	Email    string         `json:"email"`
	Sender   string         `json:"sender"`
	Filename string         `json:"filename"`
}

type Response struct {
	Sent bool `json:"sent"`
}

const TEMPLATE = `
<h1 style="color:darkgreen;">Transaction Summary</h1>
<br/>
<div>
	<ul>
		<li><strong>Total Balance:</strong> {{.TotalBalance}}</li>
		<li><strong>Average Debit:</strong> {{.AverageDebitAmount}}</li>
		<li><strong>Average Credit:</strong> {{.AverageCreditAmount}}</li>
	</ul>
	<h3>Monthly Summary:</h3>
	<ul>
	{{range .Transactions}}
		<li><strong>{{.Month}}:</strong> {{.NumberOfTransactions}} transactions</li>
	{{end}}
	</ul>
</div>
<img src="https://www.storicard.com/_next/static/media/stori_s_color.90dc745f.svg" alt="Stori Logo"/>
`

func GetFileFromS3(filename string) (io.ReadCloser, error) {
	client := s3.New(session.Must(session.NewSession()))
	response, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("stori-challenge-david-s"), //Avoid changing this bucket name pls.
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func renderHTML(summary AccountSummary) ([]byte, error) {
	tmpl, err := template.New("summaryEmail").Parse(TEMPLATE)
	if err != nil {
		return nil, err
	}
	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, summary); err != nil {
		return nil, err
	}
	return htmlBody.Bytes(), nil
}

func createAttachmentPart(writer *multipart.Writer, fileContent []byte) error {
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", "text/csv; name=\"transactions.csv\"")
	header.Set("Content-Disposition", "attachment; filename=\"transactions.csv\"")
	header.Set("Content-Transfer-Encoding", "base64")

	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(fileContent)))
	base64.StdEncoding.Encode(encoded, fileContent)
	_, err = part.Write(encoded)
	return err
}

func handler(ctx context.Context, input Request) (Response, error) {
	sess := session.Must(session.NewSession())
	emailSession := ses.New(sess)
	csvData, getError := GetFileFromS3(input.Filename)
	if getError != nil {
		return Response{false}, fmt.Errorf("unable to get file from S3: %s", getError.Error())
	}
	defer csvData.Close()

	var csvBuffer bytes.Buffer
	if _, err := io.Copy(&csvBuffer, csvData); err != nil {
		return Response{Sent: false}, fmt.Errorf("failed to read S3 object: %w", err)
	}

	htmlBody, err := renderHTML(input.Summary)
	if err != nil {
		return Response{Sent: false}, err
	}

	var rawMsg bytes.Buffer
	writer := multipart.NewWriter(&rawMsg)
	boundary := writer.Boundary()

	email := &bytes.Buffer{}
	email.WriteString(fmt.Sprintf("From: %s\n", input.Sender))
	email.WriteString(fmt.Sprintf("To: %s\n", input.Email))
	email.WriteString("Subject: Transaction Summary\n")
	email.WriteString("MIME-Version: 1.0\n")
	email.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n\n", boundary))

	mimeWriter := multipart.NewWriter(email)
	mimeWriter.SetBoundary(boundary)

	htmlHeader := textproto.MIMEHeader{}
	htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
	htmlHeader.Set("Content-Transfer-Encoding", "quoted-printable")
	htmlPart, _ := mimeWriter.CreatePart(htmlHeader)

	qp := quotedprintable.NewWriter(htmlPart)
	qp.Write(htmlBody)
	qp.Close()

	if err := createAttachmentPart(mimeWriter, csvBuffer.Bytes()); err != nil {
		return Response{Sent: false}, fmt.Errorf("failed to attach file: %w", err)
	}
	mimeWriter.Close()
	_, err = emailSession.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{Data: email.Bytes()},
	})
	if err != nil {
		return Response{Sent: false}, fmt.Errorf("failed to send email: %w", err)
	}
	return Response{Sent: true}, nil
}

func main() {
	lambda.Start(handler)
}
