package utils

import (
	"log"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
)

/* create dns resource record */
func (objMgr *ObjectManager) CreateDnsRR(payload *en.IPCDnsRRPost) error {

	_, err := objMgr.connector.CreateObject(payload, "ipcaddrr")

	if err != nil {
		log.Printf("[DEBUG] Error when create resource records , err: '%s'\n", err)
		return err
	}

	return nil
}

/* get resource record */
func (objMgr *ObjectManager) GetDnsRR(query map[string]string) (*en.IPCDnsRR, error) {
	rr := &en.IPCDnsRR{}
	queryParams := en.NewQueryParams(query)
	err := objMgr.connector.GetObject(nil, "ipcgetrr", &rr, queryParams)
	return rr, err
}

/* update dns resource record */
func (objMgr *ObjectManager) UpdateDnsRR(payload *en.IPCDnsRRPost) error {

	_, err := objMgr.connector.UpdateObject(payload, "ipcmodifyrr")

	if err != nil {
		log.Printf("[DEBUG] Error when update resource records , err: '%s'\n", err)
		return err
	}

	return nil
}

/* delete address by ip address ref */
func (objMgr *ObjectManager) DeleteDnsRR(
	owner string,
	domain string,
	resourceRecType string,
	data string,
) error {
	sf := map[string]string{
		"owner":           owner,
		"resourceRecType": resourceRecType,
		"data":            data,
		"domain":          domain,
	}
	query := en.NewQueryParams(sf)
	_, err := objMgr.connector.DeleteObject(nil, "ipcdeleterr", query)
	log.Printf("Delete DNS RR Address %s %s %s %s", owner, domain, resourceRecType, data)

	return err
}
