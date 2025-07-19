package shared

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)



func ParseCSVFromS3(event events.S3Event) ([]User, error) {
	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)

	record := event.Records[0]
	bucket := record.S3.Bucket.Name
	key := record.S3.Object.Key

	if !strings.HasPrefix(key, "uploads/") {
		return nil, fmt.Errorf("object key does not start with 'uploads/': %s", key)
	}

	resp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %v", err)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = -1

	_, _ = reader.Read() // Skip header

	var users []User
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %v", err)
		}
		if len(record) < 7 {
			return nil, fmt.Errorf("incomplete record: %v", record)
		}

		monthlyIncome, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("invalid monthly income: %v", err)
		}

		creditScore, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, fmt.Errorf("invalid credit score: %v", err)
		}

		age, err := strconv.Atoi(record[6])
		if err != nil {
			return nil, fmt.Errorf("invalid age: %v", err)
		}

		users = append(users, User{
			UserID:           record[0],
			Name:             record[1],
			Email:            record[2],
			MonthlyIncome:    monthlyIncome,
			CreditScore:      creditScore,
			EmploymentStatus: record[5],
			Age:              age,
		})
	}

	return users, nil
}