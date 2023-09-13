package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func UpdateAndArchiveItem(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse request body
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

	// Define the key to identify the item to update and archive (e.g., based on a unique ID)
	key := map[string]*dynamodb.AttributeValue{
		"field1": {S: aws.String(input.Field1)},
	}

	// Create the update expression and attribute values for the update operation
	updateExpression := "SET field2 = :field2, field3 = :field3"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":field2": {S: aws.String(input.Field2)},
		":field3": {N: aws.String(fmt.Sprintf("%d", input.Field3))},
	}

	// Create the update input for the update operation
	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	// Perform the update and archive operations
	_, err = svc.UpdateItem(updateInput)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Item updated and archived successfully",
	}, nil
}

func main() {
	lambda.Start(UpdateAndArchiveItem)
}
