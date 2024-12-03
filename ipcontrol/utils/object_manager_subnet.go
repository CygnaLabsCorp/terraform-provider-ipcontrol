package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* CreateSubnet */
func (objMgr *ObjectManager) CreateSubnet(subnet *en.IPCSubnetPost) (*en.IPCSubnetPost, error) {
	var id string
	resp, err := objMgr.connector.CreateObject(subnet, "ipcaddsubnet")
	if !strings.HasPrefix(resp, "\"") {
		resp = strconv.Quote(resp)
	}
	b := []byte(resp)

	err = json.Unmarshal(b, &id)
	if err != nil {
		log.Printf("Create Subnet Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
	}
	log.Println("[DEBUG] Subnet ID: " + fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}

	return subnet, err
}

/* get Subnet by Id ref */
func (objMgr *ObjectManager) GetSubnet(query map[string]string) (*en.IPCSubnet, error) {
	subnet := &en.IPCSubnet{}
	queryParams := en.NewQueryParams(query)
	err := objMgr.connector.GetObject(en.NewSubnet(en.IPCSubnet{}), "ipcgetsubnet", &subnet, queryParams)
	return subnet, err
}

/* delete Subnet by Id ref */
func (objMgr *ObjectManager) DeleteSubnetByIdRef(address string, size string) (string, error) {
	sf := map[string]string{
		"size":      size,
		"blockAddr": address,
	}
	query := en.NewQueryParams(sf)
	str, err := objMgr.connector.DeleteObject(en.NewSubnet(en.IPCSubnet{}), "ipcdeletechildblock", query)
	log.Printf("[DEBUG] delete subnet %s", address)
	return str, err
}

/* UpdateSubnet */
func (objMgr *ObjectManager) UpdateSubnet(
	address string,
	name string,
	size int,
	cloudType string,
	cloudObjectId string,
) (*en.IPCSubnetPost, error) {
	subnet := en.NewSubnetPost(en.IPCSubnetPost{
		Address:       address,
		Name:          name,
		Size:          size,
		CloudType:     cloudType,
		CloudObjectId: cloudObjectId,
	})

	_, err := objMgr.connector.UpdateObject(subnet, "ipcmodifysubnet")
	if err != nil {
		return nil, err
	}

	return subnet, nil
}
