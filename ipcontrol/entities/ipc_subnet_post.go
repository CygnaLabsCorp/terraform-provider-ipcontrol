package entities

type IPC_Subnet_Post struct {
	ObjBase
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Container string `json:"container,omitempty"`
	Address   string `json:"address,omitempty"`
	Type      string `json:"type,omitempty"`
	Size      string `json:"size,omitempty"`
}

/*
 * Subnet object constructor
 */
func NewSubnetPost(sb IPC_Subnet_Post) *IPC_Subnet_Post {
	res := sb
	res.objectType = "subnet"
	return &res
}

// func (b Bool) MarshalJSON() ([]byte, error) {
// 	if b {
// 		return json.Marshal("True")
// 	}
// 	return json.Marshal("False")
// }
