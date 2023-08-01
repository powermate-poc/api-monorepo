package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Input struct {
	Measurements []struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	} `json:"measurements"`
}

type SQSMessage struct {
	Message            Input  `json:"message"`
	ClientId           string `json:"clientId"`
	Timestamp_Received int64  `json:"timestamp_received"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "POST" {

		queueUrl := "https://sqs.eu-central-1.amazonaws.com/275127660632/powermate_dev_ingress_endpoint_sqs"

		device := request.PathParameters["deviceId"]

		var input Input

		err := json.Unmarshal([]byte(request.Body), &input)
		if err != nil {
			log.Printf("Error: %s", err)
			return events.APIGatewayProxyResponse{Body: "the sent body was not as expected. Bad request", StatusCode: 400}, err
		}

		sqsMessageBody := SQSMessage{Message: input, ClientId: device, Timestamp_Received: request.RequestContext.RequestTimeEpoch}

		stringified, err := json.Marshal(sqsMessageBody)
		if err != nil {
			log.Printf("Error: %s", err)
			return events.APIGatewayProxyResponse{Body: "error marshalling message", StatusCode: 500}, err
		}
		log.Println(fmt.Sprintf("sending %s to sqs", stringified))

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		svc := sqs.New(sess)
		_, err = svc.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(0),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"Author": {
					DataType:    aws.String("String"),
					StringValue: aws.String(device),
				},
			},
			MessageBody: aws.String(string(stringified)),
			QueueUrl:    &queueUrl,
		})
		if err != nil {
			log.Printf("Error: %s", err)
			return events.APIGatewayProxyResponse{Body: "error sending message to SQS", StatusCode: 500}, err
		}

		log.Println("sending to sqs successfull")

		return events.APIGatewayProxyResponse{StatusCode: 204}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 403}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
