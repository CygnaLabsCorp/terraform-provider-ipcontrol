package utils

import (
	"fmt"
	"log"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* CreateSubnet */
func (objMgr *ObjectManager) CreateQipIPv4Subnet(subnet *en.QipIPv4Subnet) (*en.QipIPv4Subnet, error) {
	_, err := objMgr.connector.CreateObject(subnet, "qipaddsubnet")
	if err != nil {
		return nil, err
	}

	return subnet, err
}

/* get Subnet by Id ref */
func (objMgr *ObjectManager) GetQipIPv4Subnet(query map[string]string) (*en.QipIPv4Subnet, error) {
	subnet := &en.QipIPv4Subnet{}
	queryParams := en.NewQueryParams(query)
	err := objMgr.connector.GetObject(en.NewQipIPv4Subnet(en.QipIPv4Subnet{}), "qipgetsubnet", &subnet, queryParams)
	log.Printf("[DEBUG] Get QIP Ipv4 Subnet: %s \n", subnet)
	return subnet, err
}

/* delete Subnet by Id ref */
func (objMgr *ObjectManager) DeleteQipIPv4Subnet(query map[string]string) error {
	queryParams := en.NewQueryParams(query)
	log.Println("[DEBUG] Subnet post: " + fmt.Sprintf("%v", queryParams))
	_, err := objMgr.connector.DeleteObject(en.NewQipIPv4Subnet(en.QipIPv4Subnet{}), "qipdeletesubnet", queryParams)
	log.Printf("delete subnet %s", query["subnetAddress"])
	return err
}

/* UpdateSubnet */
func (objMgr *ObjectManager) UpdateQipIPv4Subnet(subnet *en.QipIPv4Subnet) (*en.QipIPv4Subnet, error) {

	_, err := objMgr.connector.UpdateObject(subnet, "qipmodifysubnet")
	if err != nil {
		return nil, err
	}

	return subnet, nil
}
