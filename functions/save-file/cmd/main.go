package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
)

type File struct {
	ID        uuid.UUID       `dynamodbav:"ID"`
	S3Object  events.S3Object `dynamodbav:"s3_object"`
	Hash      string          `dynamodbav:"hash"`
	CreatedAt time.Time       `dynamodbav:"created_at"`
}

func main() {
	lambda.Start(handler)
}

const (
	tableName = "image-linker-files"
)

func handler(ctx context.Context, event events.S3Event) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := dynamodb.New(sess)

	file := File{
		ID:        uuid.New(),
		S3Object:  event.Records[0].S3.Object,
		Hash:      randstr.Hex(5),
		CreatedAt: time.Now(),
	}

	obj, err := dynamodbattribute.MarshalMap(file)
	if err != nil {
		fmt.Printf("error to marshal: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      obj,
		TableName: aws.String(tableName),
	}

	_, err = client.PutItem(input)
	if err != nil {
		fmt.Printf("error to input item: %s", err)
	}
}
