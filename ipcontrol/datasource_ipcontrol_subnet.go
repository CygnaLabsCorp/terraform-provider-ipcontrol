package ipcontrol

import (

	// "regexp"

	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetsRead,
		Schema: map[string]*schema.Schema{
			"container": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rawcontainer": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address_version": {
				ForceNew: true,
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"dns_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"block_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_object_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceSubnetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	address := strings.TrimSpace(d.Get("address").(string))
	container := strings.TrimSpace(d.Get("container").(string))
	rawContainer := d.Get("rawcontainer").(bool)
	size := d.Get("size").(int)
	status := strings.TrimSpace(d.Get("block_status").(string))

	query := map[string]string{
		"address":      address,
		"container":    container,
		"rawcontainer": strconv.FormatBool(rawContainer),
		"size":         strconv.FormatInt(int64(size), 10),
		"status":       status,
	}

	response, err := objMgr.GetSubnet(query)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Getting Subnet failed",
			Detail:   fmt.Sprintf("Getting Subnet block (%s) failed : %s", address, err),
		})
		return diags
	}

	if response == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "API returns a nil/empty Subnet",
			Detail:   fmt.Sprintf("API returns a nil/empty subnet response. Getting Subnet block (%s) failed", address),
		})
		return diags
	}

	// *** always set the ID ***
	// d.SetId("xxxx")
	// we've to definitievely store an ID on the returned resouceData object pointer, and in such operations we just  may use a timestamp,
	// we do not need a contextual ID like on CRUD operations given the nature of the "export" function ...
	// d.SetId(response.ID)
	//d.SetId(strconv.Itoa(response.ID))
	flattenIPCSubnet(d, response)
	// d.SetId(response.BlockAddr)
	//setIPCSubnetDataData(d, response)
	// subnets := flattenSubnetsData(response)
	// log.Println("[DEBUG] Subnets: " + fmt.Sprintf("%v", subnets))

	// if err := d.Set("subnets", subnets); err != nil {
	// 	return diag.FromErr(err)
	// }

	log.Println("[DEBUG] Subnet Object: " + fmt.Sprintf("%v", response))

	return nil
}

func flattenIPCSubnet(d *schema.ResourceData, subnet *en.IPCSubnet) {

	log.Println("[DEBUG] Subnet domain: " + fmt.Sprintf("%v", subnet.Subnet))

	d.SetId(strconv.Itoa(subnet.ID))
	d.Set("container", subnet.Container)
	d.Set("address", subnet.BlockAddr)
	d.Set("type", subnet.BlockType)
	d.Set("size", subnet.BlockSize)
	d.Set("dns_domain", subnet.Subnet.ForwardDomains)
	d.Set("name", subnet.BlockName)
	d.Set("block_status", subnet.BlockStatus)
	d.Set("cloud_type", subnet.CloudType)
	d.Set("cloud_object_id", subnet.CloudObjectID)
}

// func setIPCSubnetDataData(d *schema.ResourceData, subnet *en.IPCSubnet) error {
// 	// Set the container field
// 	if subnet.Container != nil {
// 		d.Set("container", subnet.Container)
// 	}

// 	// Set the numstatichosts field
// 	d.Set("numstatichosts", subnet.NumStaticHosts)

// 	// Set the subnet field (subnet is a complex type, so you will need to set it field by field)
// 	subnetMap1 := make([]map[string]interface{}, 0)
// 	subnetMap := make(map[string]interface{})
// 	subnetMap["primary_dhcp_server"] = subnet.Subnet.PrimaryDHCPServer
// 	subnetMap["cascade_primary_dhcp_server"] = subnet.Subnet.CascadePrimaryDhcpServer
// 	subnetMap["dhcp_options_set"] = subnet.Subnet.DHCPOptionsSet
// 	subnetMap["default_gateway"] = subnet.Subnet.DefaultGateway
// 	subnetMap["failover_dhcp_server"] = subnet.Subnet.FailoverDHCPServer
// 	subnetMap["dhcp_policy_set"] = subnet.Subnet.DHCPPolicySet
// 	subnetMap["forward_domains"] = subnet.Subnet.ForwardDomains
// 	subnetMap["primary_wins_server"] = subnet.Subnet.PrimaryWINSServer
// 	subnetMap["reverse_domain_types"] = subnet.Subnet.ReverseDomainTypes
// 	subnetMap["subnet_client_classes"] = subnet.Subnet.SubnetClientClasses
// 	subnetMap["reverse_domains"] = subnet.Subnet.ReverseDomains
// 	subnetMap["network_link"] = subnet.Subnet.NetworkLink
// 	subnetMap["forward_domain_types"] = subnet.Subnet.ForwardDomainTypes
// 	subnetMap["dns_servers"] = subnet.Subnet.DNSServers

// 	subnetMap1 = append(subnetMap1, subnetMap)
// 	log.Println("[DEBUG] Subnet Map: " + fmt.Sprintf("%v", subnetMap))
// 	log.Println("[DEBUG] len(subnet.VLANs): %d", len(subnet.VLANs))
// 	log.Println("[DEBUG] len(subnet.AddrDetails): %d", len(subnet.AddrDetails))

// 	d.Set("subnet", subnetMap1)

// 	// Set the VLANs (list of complex objects)
// 	if len(subnet.VLANs) > 0 {
// 		vlans := make([]interface{}, len(subnet.VLANs))
// 		for i, vlan := range subnet.VLANs {
// 			vlanMap := make(map[string]interface{})
// 			vlanMap["vlan"] = vlan.VLAN
// 			vlanMap["realm_name"] = vlan.RealmName
// 			vlanMap["vlan_name"] = vlan.VLANName
// 			vlans[i] = vlanMap
// 		}
// 		d.Set("vlans", vlans)
// 	}

// 	// Set other fields
// 	d.Set("interface_address", subnet.InterfaceAddress)
// 	d.Set("block_name", subnet.BlockName)
// 	d.Set("inherit_discovery_agent", subnet.InheritDiscoveryAgent)
// 	d.Set("num_addressable_hosts", subnet.NumAddressableHosts)
// 	d.Set("create_date", subnet.CreateDate)
// 	d.Set("description", subnet.Description)
// 	d.Set("subnet_loss_hosts", subnet.SubnetLossHosts)
// 	d.Set("num_leasable_hosts", subnet.NumLeasableHosts)
// 	d.Set("organization_id", subnet.OrganizationID)
// 	d.Set("allocation_reason_description", subnet.AllocationReasonDesc)
// 	d.Set("num_assigned_hosts", subnet.NumAssignedHosts)
// 	d.Set("root_block_type", subnet.RootBlockType)
// 	d.Set("cloud_type", subnet.CloudType)
// 	d.Set("ipv6", subnet.IPv6)
// 	d.Set("num_locked_hosts", subnet.NumLockedHosts)
// 	d.Set("rir", subnet.RIR)
// 	d.Set("interface_name", subnet.InterfaceName)
// 	d.Set("num_unallocated_hosts", subnet.NumUnallocatedHosts)
// 	d.Set("cloud_object_id", subnet.CloudObjectID)
// 	d.Set("allocation_reason", subnet.AllocationReason)
// 	d.Set("swip_name", subnet.SWIPName)
// 	d.Set("num_dynamic_hosts", subnet.NumDynamicHosts)
// 	d.Set("allow_overlapping_space", subnet.AllowOverlappingSpace)
// 	d.Set("block_type", subnet.BlockType)
// 	d.Set("discovery_agent", subnet.DiscoveryAgent)
// 	d.Set("allocation_template_name", subnet.AllocationTemplateName)
// 	d.Set("block_addr", subnet.BlockAddr)
// 	d.Set("block_status", subnet.BlockStatus)
// 	d.Set("num_reserved_hosts", subnet.NumReservedHosts)
// 	d.Set("ignore_errors", subnet.IgnoreErrors)
// 	d.Set("block_size", subnet.BlockSize)
// 	d.Set("exclude_from_discovery", subnet.ExcludeFromDiscovery)
// 	d.Set("num_allocated_hosts", subnet.NumAllocatedHosts)
// 	d.Set("user_defined_fields", subnet.UserDefinedFields)
// 	d.Set("last_admin_id", subnet.LastAdminID)
// 	d.Set("non_broadcast", subnet.NonBroadcast)
// 	d.Set("primary_subnet", subnet.PrimarySubnet)
// 	d.Set("create_admin", subnet.CreateAdmin)
// 	d.Set("keep_logical_network_link", subnet.KeepLogicalNetworkLink)

// 	// Set addr_details (list of complex objects)
// 	if len(subnet.AddrDetails) > 0 {
// 		addrDetails := make([]interface{}, len(subnet.AddrDetails))
// 		for i, addrDetail := range subnet.AddrDetails {
// 			addrDetailMap := make(map[string]interface{})
// 			addrDetailMap["starting_offset"] = addrDetail.StartingOffset
// 			addrDetailMap["share_name"] = addrDetail.ShareName
// 			addrDetailMap["offset_from_beginning_of_subnet"] = addrDetail.OffsetFromBeginningOfSubnet
// 			addrDetailMap["netservice_name"] = addrDetail.NetServiceName
// 			addrDetails[i] = addrDetailMap
// 		}
// 		d.Set("addr_details", addrDetails)
// 	}

// 	d.Set("last_update", subnet.LastUpdate)
// 	d.Set("root_block", subnet.RootBlock)

// 	return nil
// }
