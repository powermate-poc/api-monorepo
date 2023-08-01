package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
)

type QueryEvent struct {
	Query *string `json:"query"`
}

const database = "powermate-eu-central-1-dev-timestream"
const table = "powermate-eu-central-1-dev-timestream-table"
const region = "eu-central-1"

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{}

	if request.HTTPMethod != "POST" {
		return response, nil
	}

	var reqBody QueryEvent

	err := json.Unmarshal([]byte(request.Body), &reqBody)
	if err != nil {
		log.Printf("Error: %s", err)
		return events.APIGatewayProxyResponse{Body: "the sent body was not as expected. Bad request", StatusCode: 400}, err
	}

	log.Printf("Querying Timestream...")

	query, err := query(*reqBody.Query)
	if err != nil {
		log.Printf("Error: %s", err)
		return events.APIGatewayProxyResponse{Body: "error querying timestream", StatusCode: 500}, err
	}

	log.Printf("... Done")

	body, err := json.Marshal(query)

	if err != nil {
		log.Printf("Error: %s", err)
		return events.APIGatewayProxyResponse{Body: "error marshalling query response", StatusCode: 500}, err
	}

	response = events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}
	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}

func query(queryBody string) ([]*timestreamquery.Row, error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	if err != nil {
		return nil, err
	}

	// Create a new Timestream Query service client using the session
	svc := timestreamquery.New(sess)

	log.Printf("Executing Query: %s", queryBody)

	// Set up the query parameters
	params := &timestreamquery.QueryInput{
		QueryString: aws.String(queryBody)}

	// Execute the query
	resp, err := svc.Query(params)
	if err != nil {
		return nil, err
	}

	// Print the query results
	return resp.Rows, nil
}
