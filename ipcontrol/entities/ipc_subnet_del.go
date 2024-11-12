package entities

type IPCSubnetDel struct {
	ObjBase
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Name      string `json:"name,omitempty"`
	Container string `json:"container,omitempty"`
}

/*
 * Subnet object constructor
 */
func NewSubnetDel(sb IPCSubnetDel) *IPCSubnetDel {
	res := sb
	res.objectType = "subnet"
	return &res
}
