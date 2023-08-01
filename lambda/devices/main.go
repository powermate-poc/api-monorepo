package main

import (
	"encoding/json"
	"powermate-api/service/thing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	client := thing.NewClient()

	if request.HTTPMethod == "GET" {
		names, err := client.ListAll()
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error Listing Devices"}, err
		}

		body, err := json.Marshal(names)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error Marshalling Response"}, err
		}

		return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 403}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
