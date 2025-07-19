package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"www.github.com/prashantSj789/shared"
)


func handler(ctx context.Context, event events.S3Event) error {
	users, err := shared.ParseCSVFromS3(event)
	if err != nil {
		return err
	}

	db, err := shared.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, user := range users {
		if err := shared.InsertUser(db, user); err != nil {
			return err
		}
	}

	err = shared.TriggerN8NWebhook(users)
	return err

}

func presignHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    bucket := os.Getenv("UPLOAD_BUCKET")
    if bucket == "" {
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Missing UPLOAD_BUCKET env var"}, nil
    }

    key := fmt.Sprintf("uploads/%d.csv", time.Now().Unix())

    sess := session.Must(session.NewSession())
    svc := s3.New(sess)

    // Step 1: Create the PutObject request
    reqInput := &s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    }

    // Step 2: Create a request object and presign it
    req_, _ := svc.PutObjectRequest(reqInput)
    urlStr, err := req_.Presign(15 * time.Minute)
    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
    }

    body, _ := json.Marshal(map[string]string{
        "upload_url": urlStr,
        "key":        key,
    })

    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       string(body),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }, nil
}

func main() {

	lambda.Start(handler)
}