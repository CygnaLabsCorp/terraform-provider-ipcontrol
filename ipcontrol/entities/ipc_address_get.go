package entities

type IPCAddressGet struct {
	ObjBase
	AddressType string      `json:"addressType"`
	Aliases     []string    `json:"aliases"`
	Container   string      `json:"container"`
	Description string      `json:"description"`
	DeviceType  string      `json:"deviceType"`
	DUID        string      `json:"duid"`
	Hostname    string      `json:"hostname"`
	ID          int         `json:"id"`
	Interfaces  []Interface `json:"interfaces"`
	IPAddress   string      `json:"ipAddress"`
}

type Interface struct {
	AddressType          []string `json:"addressType"`
	Container            []string `json:"container"`
	ExcludeFromDiscovery string   `json:"excludeFromDiscovery"`
	ID                   int      `json:"id"`
	IPAddress            []string `json:"ipAddress"`
	Manufacturer         string   `json:"manufacturer"`
	Name                 string   `json:"name"`
	RelayAgentCircuitId  []string `json:"relayAgentCircuitId"`
	RelayAgentRemoteId   []string `json:"relayAgentRemoteId"`
	Sequence             int      `json:"sequence"`
	Virtual              []bool   `json:"virtual"`
}

func NewIPCAddressGet(sb IPCAddressGet) *IPCAddressGet {
	res := sb
	res.objectType = "address"
	return &res
}

// func (m IPCAddressGet) String() string {
// 	return fmt.Sprintf(
// 	)
// }
