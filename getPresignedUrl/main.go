package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type presignResponse struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
}

func presignHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucket := os.Getenv("UPLOAD_BUCKET")
	if bucket == "" {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Missing UPLOAD_BUCKET env var"}, nil
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	log.Println("Lambda invoked")
	log.Println("Bucket name:", bucket)

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewEnvCredentials(), // Optional: only needed outside AWS infra
	})
	if err != nil {
		log.Println("Failed to create session:", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Session error: %v", err)}, nil
	}

	svc := s3.New(sess)

	// Generate unique key
	key := fmt.Sprintf("uploads/%d.csv", time.Now().Unix())

	// Build request
	reqInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Generate presigned PUT URL (valid for 15 minutes)
	req_, _ := svc.PutObjectRequest(reqInput)
	urlStr, err := req_.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to presign request:", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Presign error: %v", err)}, nil
	}

	log.Println("Generated URL:", urlStr)

	// Return JSON response
	resp := presignResponse{
		UploadURL: urlStr,
		Key:       key,
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to serialize response"}, nil
	}



	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       strings.TrimSpace(buf.String()),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}

func main() {
	lambda.Start(presignHandler)
}

