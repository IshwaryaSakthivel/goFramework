package main

import (
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

func DeleteItem(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	field1 := request.PathParameters["field1"]

	// Create a new DynamoDB session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Create a DynamoDB client
	svc := dynamodb.New(sess)

	// Define the key to identify the item to delete based on 'field1'
	key := map[string]*dynamodb.AttributeValue{
		"field1": {S: aws.String(field1)},
	}

	// Create the delete input
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	// Perform the delete
	_, err = svc.DeleteItem(deleteInput)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Item deleted successfully",
	}, nil
}

func main() {
	lambda.Start(DeleteItem)
}
