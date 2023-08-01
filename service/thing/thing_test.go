package thing

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iot/iotiface"
	"github.com/stretchr/testify/assert"
)

type mockedClient struct {
	iotiface.IoTAPI
}

const (
	thingNameExists        = "exists"
	thingNameNotExists     = "not-exists"
	thingNameCreationError = "thing-name-creation-error"
)

func (mc mockedClient) DescribeThing(in *iot.DescribeThingInput) (*iot.DescribeThingOutput, error) {

	if *in.ThingName == thingNameNotExists {
		return nil, awserr.New("ResourceNotFoundException", "not exists", nil)
	}

	return nil, nil
}

func (mc mockedClient) CreateThing(in *iot.CreateThingInput) (*iot.CreateThingOutput, error) {
	if *in.ThingName == thingNameCreationError {
		return nil, errors.New("creation error")
	}

	return nil, nil
}

func (mc mockedClient) CreateKeysAndCertificate(in *iot.CreateKeysAndCertificateInput) (*iot.CreateKeysAndCertificateOutput, error) {

	return &iot.CreateKeysAndCertificateOutput{CertificateArn: aws.String("arn"), CertificateId: aws.String("id"), CertificatePem: aws.String("pem"), KeyPair: &iot.KeyPair{PrivateKey: aws.String("private"), PublicKey: aws.String("public")}}, nil
}

func (mc mockedClient) AttachPolicy(in *iot.AttachPolicyInput) (*iot.AttachPolicyOutput, error) {
	return nil, nil
}

func (mc mockedClient) AttachThingPrincipal(in *iot.AttachThingPrincipalInput) (*iot.AttachThingPrincipalOutput, error) {

	return nil, nil
}

func TestCreateIfNotExists(t *testing.T) {

	client := &Client{mockedClient{}}
	thingClient := ThingClient{client, "test-thing-policy-name"}

	actual, _ := thingClient.Create(thingNameNotExists)

	expected := &ThingCreationSuccess{Name: thingNameNotExists, ARN: "arn", PEM: "pem", PublicKey: "public", PrivateKey: "private", RootCA: "root-ca"}

	assert.Equal(t, actual.ARN, expected.ARN)
	assert.Equal(t, actual.PEM, expected.PEM)
	assert.Equal(t, actual.PrivateKey, expected.PrivateKey)
	assert.Equal(t, actual.PublicKey, expected.PublicKey)
}

func TestCreateExists(t *testing.T) {

	client := &Client{mockedClient{}}
	thingClient := ThingClient{client, "test-thing-policy-name"}

	_, err := thingClient.Create(thingNameExists)

	assert.Equal(t, err, ThingAlreadyExists{})
}
