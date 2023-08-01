package thing

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
)

type ThingCreationSuccess struct {
	Name       string `json:"name"`
	ARN        string `json:"arn"`
	PEM        string `json:"pem"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	RootCA     string `json:"root_ca"`
}

const ROOT_CA_URL = "https://www.amazontrust.com/repository/AmazonRootCA1.pem"

func (svc *ThingClient) Create(thingName string) (*ThingCreationSuccess, error) {

	exists, err := svc.CheckIfThingExists(thingName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ThingAlreadyExists{}
	}

	// Create a new thing
	createThingInput := &iot.CreateThingInput{
		ThingName: aws.String(thingName),
	}
	_, err = svc.core.CreateThing(createThingInput)
	if err != nil {
		log.Println("Failed to create thing:", err)
		return nil, err
	}
	log.Println("Thing created successfully")

	cert, err := svc.GenerateCertificateForThing()
	if err != nil {
		return nil, err
	}

	err = svc.AttachPolicyToCertificate(svc.ThingPolicyName, cert.ARN)
	if err != nil {
		return nil, err
	}

	err = svc.AttachCertificateWithThing(thingName, cert.ARN)

	resp, err := http.Get(ROOT_CA_URL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	root_ca := string(body)

	return &ThingCreationSuccess{Name: thingName, ARN: cert.ARN, PEM: cert.PEM, PrivateKey: cert.PrivateKey, PublicKey: cert.PublicKey, RootCA: root_ca}, nil
}

type CertificateGenerationSuccess struct {
	ARN        string
	PEM        string
	PrivateKey string
	PublicKey  string
}

func (svc *Client) GenerateCertificateForThing() (*CertificateGenerationSuccess, error) {
	createKeysAndCertificateInput := &iot.CreateKeysAndCertificateInput{
		SetAsActive: aws.Bool(true),
	}
	cert, err := svc.core.CreateKeysAndCertificate(createKeysAndCertificateInput)
	if err != nil {
		log.Println("Failed to generate certificate:", err)
		return nil, err
	}
	log.Println("Certificate generated successfully")

	return &CertificateGenerationSuccess{ARN: *cert.CertificateArn, PEM: *cert.CertificatePem, PrivateKey: *cert.KeyPair.PrivateKey, PublicKey: *cert.KeyPair.PublicKey}, nil
}

func (svc *Client) AttachPolicyToCertificate(policyName string, certificateArn string) error {
	attachPolicyInput := &iot.AttachPolicyInput{
		PolicyName: aws.String(policyName),
		Target:     aws.String(certificateArn),
	}
	_, err := svc.core.AttachPolicy(attachPolicyInput)
	if err != nil {
		log.Println("Failed to attach policy:", err)
		return err
	}
	log.Println("Policy attached successfully")

	return nil
}

func (svc *Client) AttachCertificateWithThing(thingName string, certificateArn string) error {
	attachThingPrincipalInput := &iot.AttachThingPrincipalInput{
		ThingName: aws.String(thingName),
		Principal: aws.String(certificateArn),
	}
	_, err := svc.core.AttachThingPrincipal(attachThingPrincipalInput)
	if err != nil {
		log.Println("Failed to associate certificate with the thing:", err)
		return err
	}
	log.Println("Certificate linked to the thing successfully")
	return nil
}
