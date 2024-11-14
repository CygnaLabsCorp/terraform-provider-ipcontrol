package entities

import "fmt"

type IPCSubnetPost struct {
	ObjBase        `json:"-"`
	Container      string `json:"container,omitempty"`
	RawContainer   bool   `json:"rawcontainer,omitempty"`
	Address        string `json:"address,omitempty"`
	AddressVersion int    `json:"addressversion,omitempty"`
	Type           string `json:"type,omitempty"`
	Size           int    `json:"size,omitempty"`
	DNSDomain      string `json:"dnsdomain,omitempty"`
	Name           string `json:"name,omitempty"`
	BlockStatus    string `json:"blockStatus,omitempty"`
	CloudType      string `json:"cloudType,omitempty"`
	CloudObjectId  string `json:"cloudObjectId,omitempty"`
}

/*
 * Subnet object constructor
 */
func NewSubnetPost(sb IPCSubnetPost) *IPCSubnetPost {
	res := sb
	res.objectType = "subnet"
	return &res
}

func (m IPCSubnetPost) String() string {
	return fmt.Sprintf(
		"IPCSubnetPost: Container: %s\nRawContainer: %t\nAddress: %s\nAddressVersion: %d\nType: %s\nSize: %d\nDNSDomain: %s\nName: %s\nBlockStatus: %s\nCloudType: %s\nCloudObjectId: %s",
		m.Container,
		m.RawContainer,
		m.Address,
		m.AddressVersion,
		m.Type,
		m.Size,
		m.DNSDomain,
		m.Name,
		m.BlockStatus,
		m.CloudType,
		m.CloudObjectId,
	)
}
