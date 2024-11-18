package utils

import (
	"fmt"
	"log"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* CreateSubnet */
func (objMgr *ObjectManager) CreateSubnet(subnet *en.IPCSubnetPost) (*en.IPCSubnetPost, error) {

	idRef, err := objMgr.connector.CreateObject(subnet, "ipcaddsubnet")
	log.Println("[DEBUG] Subnet ID: " + fmt.Sprintf("%v", idRef))
	if err != nil {
		return nil, err
	}

	return subnet, err
}

/* get Subnet by Id ref */
func (objMgr *ObjectManager) GetSubnet(query map[string]string) (*en.IPCSubnet, error) {
	subnet := &en.IPCSubnet{}
	queryParams := en.NewQueryParams(query)
	err := objMgr.connector.GetObject(nil, "ipcgetsubnet", &subnet, queryParams)
	return subnet, err
}

/* delete Subnet by Id ref */
func (objMgr *ObjectManager) DeleteSubnetByIdRef(address string, size string) (string, error) {
	sf := map[string]string{
		"size":      size,
		"blockAddr": address,
	}
	query := en.NewQueryParams(sf)
	str, err := objMgr.connector.DeleteObject(nil, "ipcdeletechildblock", query)
	log.Printf("delete subnet %s", address)
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
