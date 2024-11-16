package utils

import (
	"fmt"
	"log"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* CreateSubnet */
func (objMgr *ObjectManager) CreateQipIpv6Subnet(subnet *en.QipIPv6Subnet) error {
	_, err := objMgr.connector.CreateObject(subnet, "qipaddsubnet")
	if err != nil {
		return err
	}

	return nil
}

/* get Subnet by Id ref */
func (objMgr *ObjectManager) GetQipIPv6Subnet(query map[string]string) (*en.QipIPv6SubnetGet, error) {
	subnet := &en.QipIPv6SubnetGet{}
	queryParams := en.NewQueryParams(query)
	err := objMgr.connector.GetObject(nil, "qipgetsubnet", &subnet, queryParams)
	log.Printf("[DEBUG] Get QIP Ipv6 Subnet: %s \n", subnet)
	return subnet, err
}

/* delete Subnet by Id ref */
func (objMgr *ObjectManager) DeleteQipIPv6Subnet(query map[string]string) error {
	queryParams := en.NewQueryParams(query)
	log.Println("[DEBUG] Subnet delete: " + fmt.Sprintf("%v", queryParams))
	_, err := objMgr.connector.DeleteObject(nil, "qipdeletesubnet", queryParams)
	return err
}

/* UpdateSubnet */
func (objMgr *ObjectManager) UpdateQipIPv6Subnet(subnet *en.QipIPv6SubnetModify) (*en.QipIPv6SubnetModify, error) {

	_, err := objMgr.connector.UpdateObject(subnet, "qipmodifysubnet")
	if err != nil {
		return nil, err
	}

	return subnet, nil
}
