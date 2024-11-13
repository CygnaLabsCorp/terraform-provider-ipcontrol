package utils

import (
	"fmt"
	"log"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* CreateSubnet */
func (objMgr *ObjectManager) CreateSubnet(
	container string,
	address string,
	typeSubnet string,
	size int,
	name string,
	version int,
) (*en.IPCSubnetPost, error) {
	subnet := en.NewSubnetPost(en.IPCSubnetPost{
		Container:      container,
		Address:        address,
		Type:           typeSubnet,
		Size:           size,
		AddressVersion: version,
		Name:           name,
	})

	idRef, err := objMgr.connector.CreateObject(subnet, "/ipcaddsubnet")
	log.Println("[DEBUG] Subnet ID: " + fmt.Sprintf("%v", idRef))
	if err != nil {
		return nil, err
	}

	return subnet, err
}

/* get Subnet by Id ref */
func (objMgr *ObjectManager) GetSubnetByIdRef(idRef string) (*en.IPCSubnet, error) {
	subnet := &en.IPCSubnet{}
	sf := map[string]string{
		"address": idRef,
	}
	queryParams := en.NewQueryParams(sf)
	err := objMgr.connector.GetObject(nil, "/ipcgetsubnet", &subnet, queryParams)
	log.Printf("[DEBUG] get subnet: %s \n", subnet)
	return subnet, err
}

/* delete Subnet by Id ref */
func (objMgr *ObjectManager) DeleteSubnetByIdRef(idRef string, size string) (string, error) {
	sf := map[string]string{
		"size":      size,
		"blockAddr": idRef,
	}
	query := en.NewQueryParams(sf)
	str, err := objMgr.connector.DeleteObject(nil, "/ipcdeletechildblock", query)
	log.Printf("delete subnet %s", idRef)
	return str, err
}

/* UpdateSubnet */
func (objMgr *ObjectManager) UpdateSubnet(
	idRef string,
	name string,
	size int,
) (*en.IPCSubnetPost, error) {
	subnet := en.NewSubnetPost(en.IPCSubnetPost{
		Address: idRef,
		Name:    name,
		Size:    size,
	})

	idRef, err := objMgr.connector.UpdateObject(subnet, "/ipcmodifysubnet")
	if err != nil {
		return nil, err
	}

	return subnet, nil
}

/* Export Subnet(s) via body parameter selectors */
func (objMgr *ObjectManager) ExportSubnets(params en.Params) (*[]en.IPCSubnet, error) {
	subnets := []en.IPCSubnet{}

	// instantiate an empty subnet, so that the objectType will be picked by the build* functions in the caaclient.go
	subnet := en.NewSubnet(en.IPCSubnet{})

	// append all params to the subnet
	// if params != nil && len(params) > 0 {
	// 	subnet.Parameters = en.Params(params)
	// }
	query := en.NewQueryParams(nil)
	err := objMgr.connector.ExportObjects(subnet, &subnets, query)

	// log.Println("[DEBUG] ExportSubnets Response: "+fmt.Sprintf("%v", subnets))

	if err != nil || subnets == nil { // || len(res) == 0 {
		return nil, err
	}

	return &subnets, nil
}

func (objMgr *ObjectManager) GetIPAddress(ip string, container string) (*en.IPCAddressGet, error) {
	ipAddress := en.NewIPCAddressGet(en.IPCAddressGet{})
	sf := map[string]string{
		"iPAddress": ip,
		"container": container,
	}
	query := en.NewQueryParams(sf)
	err := objMgr.connector.GetObject(ipAddress, "/ipcgetdevice", &ipAddress, query)
	log.Printf("[DEBUG] get address: %v", ipAddress)
	return ipAddress, err
}
