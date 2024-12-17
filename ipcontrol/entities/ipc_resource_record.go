package entities

type IPCDnsRR struct {
	ObjBase
	TTL               string `json:"TTL"`
	Data              string `json:"data"`
	DeviceRecFlag     bool   `json:"deviceRecFlag"`
	Domain            string `json:"domain"`
	DomainType        string `json:"domainType"`
	ID                int    `json:"id"`
	Owner             string `json:"owner"`
	PendingDeployment bool   `json:"pendingDeployment"`
	ResourceRecClass  string `json:"resourceRecClass"`
	ResourceRecType   string `json:"resourceRecType"`
	Comment           string `json:"comment"`
}

/*
 * Resource record object constructor
 */
func NewResourceRecord(sb IPCDnsRR) *IPCDnsRR {
	res := sb
	res.objectType = "ipc_dns_rr"
	return &res
}

type IPCDnsRRPost struct {
	ObjBase
	Data              string `json:"data"`
	ResourceRecType   string `json:"resourceRecType"`
	Domain            string `json:"domain"`
	Owner             string `json:"owner"`
	DeviceRecFlag     bool   `json:"deviceRecFlag,omitempty"`
	DomainType        string `json:"domainType,omitempty"`
	TTL               string `json:"TTL,omitempty"`
	ID                int    `json:"id,omitempty"`
	Comment           string `json:"comment,omitempty"`
	ResourceRecClass  string `json:"resourceRecClass,omitempty"`
	PendingDeployment bool   `json:"pendingDeployment,omitempty"`
}

/*
 * Resource record object constructor
 */
func NewDnsRRPost(sb IPCDnsRRPost) *IPCDnsRRPost {
	res := sb
	res.objectType = "ipc_dns_rr"
	return &res
}
