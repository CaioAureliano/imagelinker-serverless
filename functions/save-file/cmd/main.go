package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type File struct {
	ID        string `dynamodb:"ID"`
	S3ID      string `dynamodb:"s3_id"`
	Hash      string `dynamodb:"hash"`
	CreatedAt string `dynamodb:"created_at"`
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.S3Event) {
	eventJson, _ := json.Marshal(event)

	println("event")
	fmt.Printf("%s", eventJson)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := dynamodb.New(sess)

	file := File{
		ID:        "12345",
		S3ID:      "test",
		Hash:      "zzzzza",
		CreatedAt: "now",
	}
	obj, err := dynamodbattribute.MarshalMap(file)
	if err != nil {
		fmt.Printf("error to marshal: %s", err)
	}

	tableName := "image-linker-files"

	input := &dynamodb.PutItemInput{
		Item:      obj,
		TableName: aws.String(tableName),
	}

	_, err = client.PutItem(input)
	if err != nil {
		fmt.Printf("error to input item: %s", err)
	}
}
