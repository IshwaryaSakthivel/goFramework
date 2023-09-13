package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	tableName = os.Getenv("TEST_TABLE")
	region    = os.Getenv("REGION")
)

type Samplestruct struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 int64  `json:"field3"`
}

func AddItem(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the request body into a Samplestruct
	var input Samplestruct
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	// Create a new DynamoDB session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Create a DynamoDB client
	svc := dynamodb.New(sess)

	// Marshal the Samplestruct into a DynamoDB item
	item, err := dynamodbattribute.MarshalMap(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Define the input parameters for the PutItem operation
	inputParams := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	// Perform the PutItem operation to add the item to DynamoDB
	_, err = svc.PutItem(inputParams)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201, // Created status code
		Body:       "Item added successfully",
	}, nil
}

func main() {
	lambda.Start(AddItem)
}
