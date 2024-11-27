package entities

import (
	"fmt"
)

type IPCSubnet struct {
	ObjBase
	Container     []string `json:"container"`
	BlockAddr     string   `json:"blockAddr"`
	BlockType     string   `json:"blockType"`
	BlockSize     int      `json:"blockSize"`
	BlockName     string   `json:"blockName"`
	BlockStatus   string   `json:"blockStatus"`
	CloudType     string   `json:"cloudType"`
	CloudObjectID string   `json:"cloudObjectId"`
}

/*
 * Subnet object constructor
 */
func NewSubnet(sb IPCSubnet) *IPCSubnet {
	res := sb
	res.objectType = "ipc_subnet"
	return &res
}

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
	res.objectType = "ipc_subnet_post"
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
