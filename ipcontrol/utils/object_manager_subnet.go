package utils

import (
	"fmt"
	"log"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* CreateSubnet */
func (objMgr *ObjectManager) CreateSubnet(container string, address string, typeSubnet string, size string) (*en.IPCSubnetPost, error) {
	subnet := en.NewSubnetPost(en.IPCSubnetPost{
		Container: container,
		Address:   address,
		Type:      typeSubnet,
		Size:      size,
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
	err := objMgr.connector.GetObject(nil, "/ipcgetsubnet", &subnet)
	log.Println("[DEBUG] get subnet: %s", subnet)
	return subnet, err
}

/* delete Subnet by Id ref */
func (objMgr *ObjectManager) DeleteSubnetByIdRef(idRef string) (string, error) {
	subnet := en.NewSubnetDel(en.IPCSubnetDel{
		Username:  "incadmin",
		Password:  "incadmin",
		Container: "InControl/phong",
		Name:      "138.0.0.0/24",
	})
	str, err := objMgr.connector.DeleteObject(subnet, "/ipcdeletesubnet")
	log.Printf("delete subnet %s", subnet.Name)
	return str, err
}

/* UpdateSubnet */
func (objMgr *ObjectManager) UpdateSubnet(idRef string, params en.Params) (*en.IPCSubnetPost, error) {
	subnet := en.NewSubnetPost(en.IPCSubnetPost{
		Username:  params["username"].(string),
		Password:  params["password"].(string),
		Container: params["container"].(string),
		Address:   params["address"].(string),
		Type:      params["type"].(string),
		Size:      params["size"].(string),
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

	err := objMgr.connector.ExportObjects(subnet, &subnets)

	// log.Println("[DEBUG] ExportSubnets Response: "+fmt.Sprintf("%v", subnets))

	if err != nil || subnets == nil { // || len(res) == 0 {
		return nil, err
	}

	return &subnets, nil
}
