package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ItemResult struct {
	S3Object events.S3Object `json:"s3_object,omitempty",dynamodbav:"s3_object"`
	Hash     string          `json:"hash,omitempty",dynamodbav:"hash"`
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"hash": {
				S: aws.String(request.PathParameters["hash"]),
			},
		},
		TableName: aws.String("image-linker-files"),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		fmt.Println("error to find item from dynamodb")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       fmt.Sprintf("Not found item with hash %s", request.PathParameters["hash"]),
		}, nil
	}

	var itemResult ItemResult
	_ = dynamodbattribute.UnmarshalMap(result.Item, &itemResult)

	objectInput := &s3.GetObjectInput{
		Bucket: aws.String("image-linker-files"),
		Key:    aws.String(itemResult.S3Object.Key),
	}

	s3svc := s3.New(sess)

	s3Result, err := s3svc.GetObject(objectInput)
	if err != nil {
		fmt.Println("error to get file from s3 bucket")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error to get file from bucket",
		}, nil
	}

	defer s3Result.Body.Close()

	body, err := ioutil.ReadAll(s3Result.Body)
	if err != nil {
		fmt.Println("error convert file")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "internal server error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}
