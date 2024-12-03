package entities

type AddressInterface struct {
	AddressType          []string      `json:"addressType"`
	Container            []string      `json:"container,omitempty"`
	IpAddress            []string      `json:"ipAddress"`
	ID                   int           `json:"id,omitempty"`
	Name                 string        `json:"name"`
	ExcludeFromDiscovery string        `json:"excludeFromDiscovery,omitempty"`
	Manufacturer         string        `json:"manufacturer,omitempty"`
	Sequence             int           `json:"sequence,omitempty"`
	RelayAgentRemoteId   []interface{} `json:"relayAgentRemoteId,omitempty"`
	Virtual              []interface{} `json:"virtual,omitempty"`
	RelayAgentCircuitId  []interface{} `json:"relayAgentCircuitId,omitempty"`
}

type IPCAddress struct {
	ObjBase            `json:"-"`
	ID                 int                `json:"id"`
	AddressType        string             `json:"addressType"`
	Alias              []string           `json:"alias"`
	Container          string             `json:"container"`
	Description        string             `json:"description"`
	DeviceType         string             `json:"deviceType"`
	DomainName         string             `json:"domainName"`
	DomainType         string             `json:"domainType"`
	Duid               string             `json:"duid"`
	HostName           string             `json:"hostname"`
	ResourceRecordFlag string             `json:"resourceRecordFlag"`
	IpAddress          string             `json:"ipAddress"`
	Interfaces         []AddressInterface `json:"interfaces"`
}

/*
 * Address object constructor
 */
func NewAddress(sb IPCAddress) *IPCAddress {
	res := sb
	res.objectType = "ipc_address"
	return &res
}

type IPCAddressPost struct {
	ObjBase    `json:"-"`
	Options    []string           `json:"options"`
	ID         int                `json:"id,omitempty"`
	DeviceType string             `json:"deviceType"`
	DomainName string             `json:"domainName"`
	HostName   string             `json:"hostname"`
	Interfaces []AddressInterface `json:"interfaces"`
}

/*
 * Address object constructor
 */
func NewAddressPost(sb IPCAddressPost) *IPCAddressPost {
	res := sb
	res.objectType = "ipc_address"
	return &res
}
