package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
)

type QueryEvent struct {
	Query *string `json:"query"`
}

var database string
var table string

const QUERY_CURRENT = `SELECT time, MAX(measure_value::double) AS consumption_rate FROM "%s"."%s" WHERE measure_name = '%s' AND device_id = '%s' GROUP BY time ORDER BY time DESC LIMIT 1`

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "GET" {
		return events.APIGatewayProxyResponse{StatusCode: 403}, nil
	}

	device := request.PathParameters["deviceId"]
	sensor := request.PathParameters["sensorId"]

	var q string
	if strings.Contains(request.Path, "current") {
		q = fmt.Sprintf(QUERY_CURRENT, database, table, sensor, device)
	} else {
		return events.APIGatewayProxyResponse{StatusCode: 501}, nil
	}

	log.Printf("Querying Timestream...")

	query, err := executeQuery(q)
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

	response := events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}
	return response, nil
}

func main() {
	//database = assertEnv("DATABASE")
	//table = assertEnv("TABLE")

	database = "powermate_dev_timestream"
	table = "powermate_dev_timestream_table"

	lambda.Start(HandleRequest)
}

func assertEnv(env string) string {
	value, present := os.LookupEnv(env)
	if !present {
		log.Fatalf("Error with ENV variable '%s' not set", env)
	}
	return value
}

func executeQuery(queryBody string) ([]*timestreamquery.Row, error) {

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
