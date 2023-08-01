package main

import (
	"encoding/json"
	"log"
	"os"
	"powermate-api/service/thing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var thingPolicyName string

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	client := thing.NewThingClient(thingPolicyName)

	device := request.PathParameters["deviceId"]

	if request.HTTPMethod == "PUT" {
		t, err := client.Create(device)
		if err != nil {
			if err, ok := err.(thing.ThingAlreadyExists); ok {
				return events.APIGatewayProxyResponse{StatusCode: 409, Body: err.Error()}, nil
			}
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error Creating Thing"}, err
		}

		body, err := json.Marshal(t)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error Marshalling Response"}, err
		}

		return events.APIGatewayProxyResponse{StatusCode: 201, Body: string(body)}, nil
	} else if request.HTTPMethod == "DELETE" {
		err := client.Delete(device)
		if err != nil {
			if err, ok := err.(thing.ThingDoesNotExist); ok {
				return events.APIGatewayProxyResponse{StatusCode: 404, Body: err.Error()}, nil
			}
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error Deleting Resource"}, err
		}

		return events.APIGatewayProxyResponse{StatusCode: 204, Body: ""}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 403}, nil
}

func main() {
	thingPolicyName = assertEnv("THING_POLICY_NAME")

	lambda.Start(HandleRequest)
}

func assertEnv(env string) string {
	value, present := os.LookupEnv(env)
	if !present {
		log.Fatalf("Error with ENV variable '%s' not set", env)
	}
	return value
}
