package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Printf("EVENT: %v", request)

	return events.APIGatewayProxyResponse{
		Body:       "Hello!",
		StatusCode: 200,
	}, nil
}
