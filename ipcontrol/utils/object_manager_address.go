package utils

import (
	"encoding/json"
	"fmt"
	"log"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* CreateSubnet */
func (objMgr *ObjectManager) CreateAddress(addr *en.IPCAddressPost) (*en.IPCAddress, error) {
	var address []en.IPCAddress

	resp, err := objMgr.connector.CreateObject(addr, "ipcadddevice")
	log.Println("[DEBUG] Address Resp: " + fmt.Sprintf("%v", resp))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(resp), &address)

	if err != nil {
		log.Printf("Create Device Address Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
	}

	return &address[0], err
}

/* get address */
func (objMgr *ObjectManager) GetAddress(query map[string]string) (*en.IPCAddress, error) {
	address := &en.IPCAddress{}
	queryParams := en.NewQueryParams(query)
	err := objMgr.connector.GetObject(nil, "ipcgetdevice", &address, queryParams)
	return address, err
}

/* delete address by ip address ref */
func (objMgr *ObjectManager) DeleteAddressRef(address string) (string, error) {
	sf := map[string]string{
		"ipAddress": address,
	}
	query := en.NewQueryParams(sf)
	str, err := objMgr.connector.DeleteObject(en.NewAddress(en.IPCAddress{}), "ipcdeletedevice", query)
	log.Printf("Delete Device Address %s", address)
	return str, err
}

/* UpdateAddress */
func (objMgr *ObjectManager) UpdateAddress(data *en.IPCAddressPost) (*en.IPCAddress, error) {
	var address []en.IPCAddress
	resp, err := objMgr.connector.UpdateObject(data, "ipcmodifydevice")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(resp), &address)

	if err != nil {
		log.Printf("Update Device Address Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
	}

	return &address[0], err
}
