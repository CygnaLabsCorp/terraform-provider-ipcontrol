package utils

import (
	"fmt"
	"log"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

type ObjectManager struct {
	connector CAAConnector
}

func NewObjectManager(connector CAAConnector) *ObjectManager {
	objMgr := new(ObjectManager)
	objMgr.connector = connector
	return objMgr
}

/* CreateSubnet */
func (objMgr *ObjectManager) CreateSubnet(cidr string, name string, gateway string, params en.Params) (*en.Subnet, error) {
	subnet := en.NewSubnet(en.Subnet{
		Cidr: cidr,
	})

	if name != "" {
		subnet.Name = name
	}

	if gateway != "" {
		subnet.Gateway = gateway
	}

	if params != nil && len(params) > 0 {
		subnet.Parameters = en.Params(params)
		log.Println("[DEBUG] PARAMS: " + fmt.Sprintf("%v", subnet.Parameters))
	}

	idRef, err := objMgr.connector.CreateObject(subnet)
	log.Println("[DEBUG] Subnet ID: " + fmt.Sprintf("%v", idRef))
	if err != nil {
		return nil, err
	}
	// sets the Id property here
	subnet.Id = idRef

	return subnet, err
}

/* get Subnet by Id ref */
func (objMgr *ObjectManager) GetSubnetByIdRef(idRef string) (*en.Subnet, error) {
	subnet := en.NewSubnet(en.Subnet{})
	err := objMgr.connector.GetObject(subnet, idRef, &subnet)
	return subnet, err
}

/* delete Subnet by Id ref */
func (objMgr *ObjectManager) DeleteSubnetByIdRef(idRef string) (string, error) {
	return objMgr.connector.DeleteObject(idRef)
}

/* UpdateSubnet */
func (objMgr *ObjectManager) UpdateSubnet(idRef string, params en.Params) (*en.Subnet, error) {
	subnet := en.NewSubnet(en.Subnet{})

	if params != nil && len(params) > 0 {
		subnet.Parameters = en.Params(params)
	}

	if params != nil && len(params) > 0 {
		subnet.Parameters = en.Params(params)
		log.Println("[DEBUG] PARAMS:: " + fmt.Sprintf("%v", subnet.Parameters))
	}

	log.Println("[DEBUG] PARAMS-UPD: " + fmt.Sprintf("%v", subnet.Parameters))
	log.Println("[DEBUG] PARAMS-UPD: " + fmt.Sprintf("%v", idRef))

	idRef, err := objMgr.connector.UpdateObject(subnet, idRef)
	if err != nil {
		return nil, err
	}
	// sets the Id property here
	subnet.Id = idRef

	return subnet, err
}

/* Export Subnet(s) via body parameter selectors */
func (objMgr *ObjectManager) ExportSubnets(params en.Params) (*[]en.Subnet, error) {
	subnets := []en.Subnet{}

	// instantiate an empty subnet, so that the objectType will be picked by the build* functions in the caaclient.go
	subnet := en.NewSubnet(en.Subnet{})

	// append all params to the subnet
	if params != nil && len(params) > 0 {
		subnet.Parameters = en.Params(params)
	}

	err := objMgr.connector.ExportObjects(subnet, &subnets)

	// log.Println("[DEBUG] ExportSubnets Response: "+fmt.Sprintf("%v", subnets))

	if err != nil || subnets == nil { // || len(res) == 0 {
		return nil, err
	}

	return &subnets, nil
}
