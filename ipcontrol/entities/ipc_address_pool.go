package entities

type IPCAddressPool struct {
	ObjBase
	Container          string `json:"container"`
	StartAddr          string `json:"startAddr"`
	EndAddr            string `json:"endAddr"`
	Type               string `json:"type"`
	Name               string `json:"name"`
	ID                 string `json:"id"`
	PrefixLength       string `json:"prefixLength"`
	CreatedDate        string `json:"createDate"`
	PrimaryNetService  string `json:"primaryNetService"`
	OverlapInterfaceIp Bool   `json:"overlapInterfaceIp"`
	LastAdmin          string `json:"lastAdmin"`
	DhcpOptionSet      string `json:"dhcpOptionSet"`
	DhcpPolicySet      string `json:"dhcpPolicySet"`
}

/*
 * Subnet object constructor
 */
func NewAddressPool(sb IPCAddressPool) *IPCAddressPool {
	res := sb
	res.objectType = "ipc_address_pool"
	return &res
}

type IPCAddressPoolPost struct {
	ObjBase            `json:"-"`
	Container          string `json:"container,omitempty"`
	StartAddr          string `json:"startAddr"`
	EndAddr            string `json:"endAddr,omitempty"`
	Type               string `json:"type"`
	ID                 string `json:"id,omitempty"`
	PrefixLength       int    `json:"prefixLength,omitempty"`
	PrimaryNetService  string `json:"primaryNetService,omitempty"`
	DhcpOptionSet      string `json:"dhcpOptionSet,omitempty"`
	DhcpPolicySet      string `json:"dhcpPolicySet,omitempty"`
	Name               string `json:"name,omitempty"`
	OverlapInterfaceIp bool   `json:"overlapInterfaceIp,omitempty"`
}

/*
 * Subnet object constructor
 */
func NewAddressPoolPost(sb IPCAddressPoolPost) *IPCAddressPoolPost {
	res := sb
	res.objectType = "ipc_address_pool_post"
	return &res
}
