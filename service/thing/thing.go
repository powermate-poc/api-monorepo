package thing

import (
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iot/iotiface"
)

type Client struct {
	core iotiface.IoTAPI
}

type ThingClient struct {
	*Client

	ThingPolicyName string
}

func NewClient() *Client {

	//os.Setenv("AWS_ACCESS_KEY_ID", "XXX")
	//os.Setenv("AWS_SECRET_ACCESS_KEY", "XXX")

	mySession := session.Must(session.NewSession())

	// Create a IoT client with additional configuration
	svc := iot.New(mySession, aws.NewConfig())

	return &Client{svc}
}

func NewThingClient(thingPolicyName string) *ThingClient {

	return &ThingClient{NewClient(), thingPolicyName}
}

type ThingAlreadyExists struct{}

type ThingDoesNotExist struct{}

func (m ThingAlreadyExists) Error() string {
	return "Thing already exists!"
}

func (m ThingDoesNotExist) Error() string {
	return "Thing does not exist!"
}

func (svc *Client) CheckIfThingExists(thingName string) (bool, error) {
	log.Printf("Checking if thing '%s' exists ...\n", thingName)
	// Check if Thing already exists (name must be unique)
	_, err := svc.core.DescribeThing(&iot.DescribeThingInput{ThingName: aws.String(thingName)})
	exists := true
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == iot.ErrCodeResourceNotFoundException {
			exists = false
		} else {
			log.Println("Error describing thing:", err)
			return true, aerr
		}
	}

	if exists {
		log.Printf("Thing '%s' does exist, aborting.\n", thingName)
		return true, nil
	} else {
		log.Printf("Thing '%s' does NOT exist, continuing...\n", thingName)
		return false, nil
	}
}

func (svc *Client) GetCertificateForThing(thingName string) (*string, *string, error) {
	out, err := svc.core.ListThingPrincipals(&iot.ListThingPrincipalsInput{ThingName: &thingName})
	if err != nil {
		log.Println("Failed to list principals for thing:", err)
		return nil, nil, err
	}
	principals := out.Principals
	if len(principals) != 1 {
		log.Println("Number of principals attached to thing not as expected")
		return nil, nil, errors.New("thing %s has none or more than one expected principal attached")
	}

	arn := *principals[0]

	id := strings.Split(arn, ":cert/")[1]

	return &arn, &id, nil
}
