package entities

import "encoding/json"

type Subnet struct {
	ObjBase
	Id          string `json:"id,omitempty"`
	Cidr        string `json:"cidr,omitempty"`
	Name        string `json:"name,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
	Blocktype   string `json:"blocktype,omitempty"`
	Blockstatus string `json:"blockstatus,omitempty"`
	Description string `json:"description,omitempty"`
	Cloudtype   string `json:"cloudtype,omitempty"`
	Cloudobjid  string `json:"cloudobjid,omitempty"`
	Tenantid    string `json:"tenantid,omitempty"`
	Tenantname  string `json:"tenantname,omitempty"`
	Container   string `json:"container,omitempty"`
	Lastupdated string `json:"lastupdated,omitempty"`
}

/*
 * Subnet object constructor
 */
func NewSubnet(sb Subnet) *Subnet {
	res := sb
	res.objectType = "subnet"
	return &res
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal("True")
	}
	return json.Marshal("False")
}
