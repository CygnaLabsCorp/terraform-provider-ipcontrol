package utils

import (
	"fmt"
	"log"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* create address pool */
func (objMgr *ObjectManager) CreateAddressPool(payload *en.IPCAddressPoolPost) (*en.IPCAddressPool, error) {
	var addressPool en.IPCAddressPool

	// can unmarshal because api return empty
	resp, err := objMgr.connector.CreateObject(payload, "/ipcimportaddrpool")

	log.Println("[DEBUG] Address Pool Create : " + fmt.Sprintf("%v", resp))

	if err != nil {
		return nil, err
	}
	// b := []byte(resp)
	// err = json.Unmarshal(b, &addressPool)

	return &addressPool, err
}

/* get address pool by start address */
func (objMgr *ObjectManager) GetAddressPool(query map[string]string) (*en.IPCAddressPool, error) {
	var addressPool en.IPCAddressPool
	queryParams := en.NewQueryParams(query)
	err := objMgr.connector.GetObject(nil, "/ipcgetaddrpool", &addressPool, queryParams)

	if err != nil {
		log.Printf("[DEBUG] Error when get address pool: %v \n", err)
		return nil, err
	}

	return &addressPool, nil
}

/* update address pool */
func (objMgr *ObjectManager) UpdateAddressPool(payload *en.IPCAddressPoolPost) (*en.IPCAddressPool, error) {
	var addressPool en.IPCAddressPool
	_, err := objMgr.connector.UpdateObject(payload, "ipcmodifyaddrpool")

	if err != nil {
		log.Printf("[DEBUG] Error when update address pool: %v \n", err)
		return nil, err
	}

	// handle resp if exist
	// there's no resp from api

	//---------------

	return &addressPool, nil
}

/* delete addresspool by start address and container */
func (objMgr *ObjectManager) DeleteAddressPool(startAddress string, container string) (string, error) {
	log.Printf("[DEBUG] Delete Device Address Pool %s \n", startAddress)

	sf := map[string]string{
		"startAddr":               startAddress,
		"instantDDNS":             "true",
		"deleteDevicesInAddrpool": "true",
		"container":               container,
	}
	query := en.NewQueryParams(sf)
	str, err := objMgr.connector.DeleteObject(nil, "ipcdeleteaddrpool", query)
	return str, err
}
