package thing

import (
	"log"

	"github.com/aws/aws-sdk-go/service/iot"
)

func (svc *ThingClient) Delete(thingName string) error {

	exists, err := svc.CheckIfThingExists(thingName)
	if err != nil {
		return err
	}
	if !exists {
		return ThingDoesNotExist{}
	}

	arn, id, err := svc.GetCertificateForThing(thingName)
	if err != nil {
		return err
	}

	err = svc.DeactivateCertificate(*id)
	if err != nil {
		return err
	}

	err = svc.DetachPolicyFromCertificate(svc.ThingPolicyName, *arn)
	if err != nil {
		return err
	}

	err = svc.DetachCertificateFromThing(*arn, thingName)
	if err != nil {
		return err
	}

	err = svc.delete(thingName, *id)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Client) DeactivateCertificate(certificateId string) error {
	inactive := "INACTIVE"
	_, err := svc.core.UpdateCertificate(&iot.UpdateCertificateInput{CertificateId: &certificateId, NewStatus: &inactive})

	if err != nil {
		log.Println("Error deactivating certficiate for thing:", err)
		return err
	}

	log.Println("Certificate deactivated successfully")

	return nil
}

func (svc *Client) DetachPolicyFromCertificate(policyName string, certificateArn string) error {
	_, err := svc.core.DetachPolicy(&iot.DetachPolicyInput{PolicyName: &policyName, Target: &certificateArn})

	if err != nil {
		log.Println("Error detaching policy from certificate:", err)
		return err
	}

	log.Println("Policy detached from certificate successfully")

	return nil
}

func (svc *Client) DetachCertificateFromThing(certificateArn string, thingName string) error {

	_, err := svc.core.DetachThingPrincipal(&iot.DetachThingPrincipalInput{Principal: &certificateArn, ThingName: &thingName})

	if err != nil {
		log.Println("Error detaching certificate from thing:", err)
		return err
	}

	log.Println("Certificate detached from thing successfully")

	return nil
}

func (svc *Client) delete(thingName string, certifcateId string) error {
	_, err := svc.core.DeleteCertificate(&iot.DeleteCertificateInput{CertificateId: &certifcateId})
	if err != nil {
		log.Println("Error deleting certificate:", err)
		return err
	}

	log.Println("Certificate deleted successfully")

	_, err = svc.core.DeleteThing(&iot.DeleteThingInput{ThingName: &thingName})
	if err != nil {
		log.Println("Error deleting thing:", err)
		return err
	}

	log.Println("Thing deleted successfully")

	return nil
}
