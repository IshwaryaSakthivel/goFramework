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

func GetItemsByField1(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get the value of field1 from the path parameters
	field1Value := request.PathParameters["field1"]

	// Create a new DynamoDB session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Create a DynamoDB client
	svc := dynamodb.New(sess)

	// Define the input parameters for the Scan operation with a filter expression
	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("field1 = :value"), // Filter by field1
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":value": {
				S: aws.String(field1Value), // Use the field1Value from path parameters
			},
		},
	}

	// Perform the Scan operation to retrieve items matching the filter
	result, err := svc.Scan(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Marshal the result items into JSON format
	items := []Samplestruct{}
	for _, item := range result.Items {
		var sample Samplestruct
		if err := dynamodbattribute.UnmarshalMap(item, &sample); err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		items = append(items, sample)
	}

	// Convert the items to JSON
	responseBody, err := json.Marshal(items)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(GetItemsByField1)
}
