package entities

type IPCSubnetPost struct {
	ObjBase        `json:"-"`
	Container      string `json:"container,omitempty"`
	Address        string `json:"address,omitempty"`
	Type           string `json:"type,omitempty"`
	Size           int    `json:"size,omitempty"`
	Name           string `json:"name,omitempty"`
	AddressVersion int    `json:"addressversion,omitempty"`
}

/*
 * Subnet object constructor
 */
func NewSubnetPost(sb IPCSubnetPost) *IPCSubnetPost {
	res := sb
	res.objectType = "subnet"
	return &res
}
