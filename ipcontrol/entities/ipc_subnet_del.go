package entities

type IPC_Subnet_Del struct {
	ObjBase
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Name      string `json:"name,omitempty"`
	Container string `json:"container,omitempty"`
}

/*
 * Subnet object constructor
 */
func NewSubnetDel(sb IPC_Subnet_Del) *IPC_Subnet_Del {
	res := sb
	res.objectType = "subnet"
	return &res
}
