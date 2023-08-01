package thing

import (
	"log"

	"github.com/aws/aws-sdk-go/service/iot"
)

func (svc *Client) ListAll() ([]Device, error) {

	things, err := svc.core.ListThings(&iot.ListThingsInput{}) //TODO: later on, check for attributes

	if err != nil {
		log.Println("Error listing things:", err)
		return nil, err
	}

	var devices []Device
	for _, v := range things.Things {
		name := *v.ThingName
		log.Println(name)
		devices = append(devices, Device{Name: name})
	}

	return devices, nil
}
