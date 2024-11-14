// package entities

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strings"
// )

// type IPCSubnet struct {
// 	ObjBase
// 	Container []string `json:"container,omitempty"`
// 	Address   string   `json:"blockAddr,omitempty"`
// 	Type      string   `json:"blockType,omitempty"`
// 	Size      int      `json:"blockSize,omitempty"`
// 	Name      string   `json:"blockName,omitempty"`
// }

// /*
//  * Subnet object constructor
//  */
// func NewSubnet(sb IPCSubnet) *IPCSubnet {
// 	res := sb
// 	res.objectType = "subnet"
// 	return &res
// }

// func (b Bool) MarshalJSON() ([]byte, error) {
// 	if b {
// 		return json.Marshal("True")
// 	}
// 	return json.Marshal("False")
// }

// func (m IPCSubnet) String() string {
// 	return fmt.Sprintf(
// 		"IPC_Subnet{Container: [%s], Address: %s, Type: %s, Size: %d}",
// 		strings.Join(m.Container, ", "),
// 		m.Address,
// 		m.Type,
// 		m.Size,
// 	)
// }

package entities

import (
	"fmt"
	"strings"
)

type IPCSubnet struct {
	ObjBase
	Container              []string     `json:"container"`
	NumStaticHosts         int          `json:"numstatichosts"`
	Subnet                 Subnet       `json:"subnet"`
	VLANs                  []VLAN       `json:"vlans"`
	InterfaceAddress       []string     `json:"interfaceAddress"`
	BlockName              string       `json:"blockName"`
	InheritDiscoveryAgent  int          `json:"inheritDiscoveryAgent"`
	NumAddressableHosts    int          `json:"numaddressablehosts"`
	CreateDate             string       `json:"createdate"`
	Description            string       `json:"description"`
	SubnetLossHosts        int          `json:"subnetlosshosts"`
	NumLeasableHosts       int          `json:"numleasablehosts"`
	OrganizationID         string       `json:"organizationId"`
	AllocationReasonDesc   string       `json:"allocationReasonDescription"`
	NumAssignedHosts       int          `json:"numassignedhosts"`
	RootBlockType          string       `json:"rootBlocktype"`
	CloudType              string       `json:"cloudType"`
	IPv6                   bool         `json:"ipv6"`
	NumLockedHosts         int          `json:"numlockedhosts"`
	RIR                    string       `json:"rir"`
	ID                     int          `json:"id"`
	InterfaceName          []string     `json:"interfaceName"`
	NumUnallocatedHosts    int          `json:"numunallocatedhosts"`
	CloudObjectID          string       `json:"cloudObjectId"`
	AllocationReason       string       `json:"allocationReason"`
	SWIPName               string       `json:"swipname"`
	NumDynamicHosts        int          `json:"numdynamichosts"`
	AllowOverlappingSpace  bool         `json:"allowOverlappingSpace"`
	BlockType              string       `json:"blockType"`
	DiscoveryAgent         string       `json:"discoveryAgent"`
	AllocationTemplateName string       `json:"allocationTemplateName"`
	BlockAddr              string       `json:"blockAddr"`
	BlockStatus            string       `json:"blockStatus"`
	NumReservedHosts       int          `json:"numreservedhosts"`
	IgnoreErrors           bool         `json:"ignoreErrors"`
	BlockSize              int          `json:"blockSize"`
	ExcludeFromDiscovery   string       `json:"excludeFromDiscovery"`
	NumAllocatedHosts      int          `json:"numallocatedhosts"`
	UserDefinedFields      []string     `json:"userDefinedFields"`
	LastAdminID            int          `json:"lastadminid"`
	NonBroadcast           bool         `json:"nonBroadcast"`
	PrimarySubnet          bool         `json:"primarySubnet"`
	CreateAdmin            string       `json:"createadmin"`
	KeepLogicalNetworkLink bool         `json:"keepLogicalNetworkLink"`
	AddrDetails            []AddrDetail `json:"addrDetails"`
	LastUpdate             string       `json:"lastupdate"`
	RootBlock              bool         `json:"rootBlock"`
}

type Subnet struct {
	PrimaryDHCPServer        string   `json:"primaryDHCPServer,omitempty"`
	CascadePrimaryDhcpServer bool     `json:"cascadePrimaryDhcpServer,omitempty"`
	DHCPOptionsSet           string   `json:"DHCPOptionsSet,omitempty"`
	DefaultGateway           string   `json:"defaultGateway,omitempty"`
	FailoverDHCPServer       string   `json:"failoverDHCPServer,omitempty"`
	DHCPPolicySet            string   `json:"DHCPPolicySet,omitempty"`
	ForwardDomains           []string `json:"forwardDomains,omitempty"`
	PrimaryWINSServer        string   `json:"primaryWINSServer,omitempty"`
	ReverseDomainTypes       []string `json:"reverseDomainTypes,omitempty"`
	SubnetClientClasses      []string `json:"subnetClientClasses,omitempty"`
	ReverseDomains           []string `json:"reverseDomains,omitempty"`
	NetworkLink              string   `json:"networkLink,omitempty"`
	ForwardDomainTypes       []string `json:"forwardDomainTypes,omitempty"`
	DNSServers               []string `json:"DNSServers,omitempty"`
}

type VLAN struct {
	VLAN      int    `json:"vlan"`
	RealmName string `json:"realmName"`
	VLANName  string `json:"vlanName"`
}

type AddrDetail struct {
	StartingOffset              int    `json:"startingOffset"`
	ShareName                   string `json:"sharename"`
	OffsetFromBeginningOfSubnet bool   `json:"offsetFromBeginningOfSubnet"`
	NetServiceName              string `json:"netserviceName"`
}

/*
 * Subnet object constructor
 */
func NewSubnet(sb IPCSubnet) *IPCSubnet {
	res := sb
	res.objectType = "subnet"
	return &res
}

func (ipc IPCSubnet) String() string {
	var sb strings.Builder

	sb.WriteString("IPCSubnet: {\n")
	sb.WriteString(fmt.Sprintf("  Container: %v\n", ipc.Container))
	sb.WriteString(fmt.Sprintf("  NumStaticHosts: %d\n", ipc.NumStaticHosts))
	sb.WriteString(fmt.Sprintf("  Subnet: %v\n", ipc.Subnet))
	sb.WriteString(fmt.Sprintf("  VLANs: %v\n", ipc.VLANs))
	sb.WriteString(fmt.Sprintf("  InterfaceAddress: %v\n", ipc.InterfaceAddress))
	sb.WriteString(fmt.Sprintf("  BlockName: %s\n", ipc.BlockName))
	sb.WriteString(fmt.Sprintf("  InheritDiscoveryAgent: %d\n", ipc.InheritDiscoveryAgent))
	sb.WriteString(fmt.Sprintf("  NumAddressableHosts: %d\n", ipc.NumAddressableHosts))
	sb.WriteString(fmt.Sprintf("  CreateDate: %s\n", ipc.CreateDate))
	sb.WriteString(fmt.Sprintf("  Description: %s\n", ipc.Description))
	sb.WriteString(fmt.Sprintf("  SubnetLossHosts: %d\n", ipc.SubnetLossHosts))
	sb.WriteString(fmt.Sprintf("  NumLeasableHosts: %d\n", ipc.NumLeasableHosts))
	sb.WriteString(fmt.Sprintf("  OrganizationID: %s\n", ipc.OrganizationID))
	sb.WriteString(fmt.Sprintf("  AllocationReasonDesc: %s\n", ipc.AllocationReasonDesc))
	sb.WriteString(fmt.Sprintf("  NumAssignedHosts: %d\n", ipc.NumAssignedHosts))
	sb.WriteString(fmt.Sprintf("  RootBlockType: %s\n", ipc.RootBlockType))
	sb.WriteString(fmt.Sprintf("  CloudType: %s\n", ipc.CloudType))
	sb.WriteString(fmt.Sprintf("  IPv6: %v\n", ipc.IPv6))
	sb.WriteString(fmt.Sprintf("  NumLockedHosts: %d\n", ipc.NumLockedHosts))
	sb.WriteString(fmt.Sprintf("  RIR: %s\n", ipc.RIR))
	sb.WriteString(fmt.Sprintf("  ID: %d\n", ipc.ID))
	sb.WriteString(fmt.Sprintf("  InterfaceName: %v\n", ipc.InterfaceName))
	sb.WriteString(fmt.Sprintf("  NumUnallocatedHosts: %d\n", ipc.NumUnallocatedHosts))
	sb.WriteString(fmt.Sprintf("  CloudObjectID: %s\n", ipc.CloudObjectID))
	sb.WriteString(fmt.Sprintf("  AllocationReason: %s\n", ipc.AllocationReason))
	sb.WriteString(fmt.Sprintf("  SWIPName: %s\n", ipc.SWIPName))
	sb.WriteString(fmt.Sprintf("  NumDynamicHosts: %d\n", ipc.NumDynamicHosts))
	sb.WriteString(fmt.Sprintf("  AllowOverlappingSpace: %v\n", ipc.AllowOverlappingSpace))
	sb.WriteString(fmt.Sprintf("  BlockType: %s\n", ipc.BlockType))
	sb.WriteString(fmt.Sprintf("  DiscoveryAgent: %s\n", ipc.DiscoveryAgent))
	sb.WriteString(fmt.Sprintf("  AllocationTemplateName: %s\n", ipc.AllocationTemplateName))
	sb.WriteString(fmt.Sprintf("  BlockAddr: %s\n", ipc.BlockAddr))
	sb.WriteString(fmt.Sprintf("  BlockStatus: %s\n", ipc.BlockStatus))
	sb.WriteString(fmt.Sprintf("  NumReservedHosts: %d\n", ipc.NumReservedHosts))
	sb.WriteString(fmt.Sprintf("  IgnoreErrors: %v\n", ipc.IgnoreErrors))
	sb.WriteString(fmt.Sprintf("  BlockSize: %d\n", ipc.BlockSize))
	sb.WriteString(fmt.Sprintf("  ExcludeFromDiscovery: %s\n", ipc.ExcludeFromDiscovery))
	sb.WriteString(fmt.Sprintf("  NumAllocatedHosts: %d\n", ipc.NumAllocatedHosts))
	sb.WriteString(fmt.Sprintf("  UserDefinedFields: %v\n", ipc.UserDefinedFields))
	sb.WriteString(fmt.Sprintf("  LastAdminID: %d\n", ipc.LastAdminID))
	sb.WriteString(fmt.Sprintf("  NonBroadcast: %v\n", ipc.NonBroadcast))
	sb.WriteString(fmt.Sprintf("  PrimarySubnet: %v\n", ipc.PrimarySubnet))
	sb.WriteString(fmt.Sprintf("  CreateAdmin: %s\n", ipc.CreateAdmin))
	sb.WriteString(fmt.Sprintf("  KeepLogicalNetworkLink: %v\n", ipc.KeepLogicalNetworkLink))
	sb.WriteString(fmt.Sprintf("  AddrDetails: %v\n", ipc.AddrDetails))
	sb.WriteString(fmt.Sprintf("  LastUpdate: %s\n", ipc.LastUpdate))
	sb.WriteString(fmt.Sprintf("  RootBlock: %v\n", ipc.RootBlock))
	sb.WriteString("}\n")

	return sb.String()
}

// String method for Subnet struct
func (s Subnet) String() string {
	return fmt.Sprintf("PrimaryDHCPServer: %s, DefaultGateway: %s, DNSServers: %v",
		s.PrimaryDHCPServer, s.DefaultGateway, s.DNSServers)
}

// String method for VLAN struct
func (v VLAN) String() string {
	return fmt.Sprintf("VLAN: %d, VLANName: %s, RealmName: %s", v.VLAN, v.VLANName, v.RealmName)
}

// String method for AddrDetail struct
func (a AddrDetail) String() string {
	return fmt.Sprintf("StartingOffset: %d, ShareName: %s, NetServiceName: %s", a.StartingOffset, a.ShareName, a.NetServiceName)
}
