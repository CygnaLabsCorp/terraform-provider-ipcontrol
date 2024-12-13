package ipcontrol

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	// "strconv"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAddressPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAddressPoolRecordContext,
		ReadContext:   getAddressPoolRecordContext,
		UpdateContext: updateAddressPoolRecordContext,
		DeleteContext: deleteAddressPoolRecordContext,

		Schema: map[string]*schema.Schema{
			"container": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the container that will hold the address pool.",
			},
			"start_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The start address of the address pool.",
			},
			"end_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The end address of the address pool.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of address pool.",
			},
			"prefix_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The size parameter represents the subnet mask or prefix length of the address block in CIDR notation. Required for IPv6, the size value is larger due to the increased address space. IPv6 prefix lengths commonly range between /48 to /128, with /64 often used as the standard size for a single subnet.",
			},
			"primary_net_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the DHCP server that will serve addresses from this pool.",
			},
			"overlap_interface_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to allow a DHCPv6 pool to overlap an interface address.",
			},
			"dhcp_policy_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a Policy Set used with this pool.",
			},
			"dhcp_option_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a Option Set used with this pool.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the address pool.",
			},
			"last_admin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest admin modified this address pool.",
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: func(ctx context.Context, rd *schema.ResourceDiff, i interface{}) error {
			startAddr := rd.Get("start_address").(string)

			if isIPv4(startAddr) {
				_, ok := rd.GetOk("end_address")
				if !ok {
					return fmt.Errorf("end_address is required for IPv4")
				}
			} else if isIPv6(startAddr) {

				if _, ok := rd.GetOk("prefix_length"); !ok {
					return fmt.Errorf("prefix_length is required for IPv6")
				}

			} else {
				return fmt.Errorf("start_address must be valid IPv4 or IPv6 address, but got %v", startAddr)
			}

			return nil
		},
	}
}

func setIPCAddressPoolResource(d *schema.ResourceData, addressPool en.IPCAddressPool) {
	d.SetId(addressPool.ID)

	d.Set("start_address", addressPool.StartAddr)
	d.Set("container", addressPool.Container)
	d.Set("type", addressPool.Type)
	d.Set("name", addressPool.Name)
	d.Set("overlap_interface_ip", addressPool.OverlapInterfaceIp)
	d.Set("last_admin", addressPool.LastAdmin)
	d.Set("create_date", addressPool.CreatedDate)
	d.Set("dhcp_policy_set", addressPool.DhcpPolicySet)
	d.Set("dhcp_option_set", addressPool.DhcpOptionSet)

	if addressPool.EndAddr != "" {
		d.Set("end_address", addressPool.EndAddr)
	}

	if addressPool.DhcpPolicySet != "" {
		d.Set("dhcp_policy_set", addressPool.DhcpPolicySet)
	}

	if addressPool.DhcpOptionSet != "" {
		d.Set("dhcp_option_set", addressPool.DhcpOptionSet)
	}

	if addressPool.PrefixLength != "" {
		prefixLength, err := strconv.Atoi(addressPool.PrefixLength)
		if err == nil {
			d.Set("prefix_length", prefixLength)
		}
	}

}

func createAddressPoolRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	startAddress := d.Get("start_address").(string)
	endAddress := d.Get("end_address").(string)
	addressPoolType := d.Get("type").(string)
	prefixLength := d.Get("prefix_length").(int)
	primaryNetService := d.Get("primary_net_service").(string)
	dhcpOptionSet := d.Get("dhcp_option_set").(string)
	dhcpPolicySet := d.Get("dhcp_policy_set").(string)
	overlapInterfaceIp := d.Get("overlap_interface_ip").(bool)
	name := d.Get("name").(string)

	payload := en.NewAddressPoolPost(en.IPCAddressPoolPost{
		StartAddr:          startAddress,
		EndAddr:            endAddress,
		Name:               name,
		Type:               addressPoolType,
		PrefixLength:       prefixLength,
		PrimaryNetService:  primaryNetService,
		DhcpOptionSet:      dhcpOptionSet,
		DhcpPolicySet:      dhcpPolicySet,
		OverlapInterfaceIp: overlapInterfaceIp,
	})

	log.Println("[DEBUG] Address pool  post: " + fmt.Sprintf("%v", payload))

	_, err := objMgr.CreateAddressPool(payload)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when create address pool",
			Detail:   fmt.Sprintf("Error when create address pool: (%s) ", err),
		})
		return diags
	}
	d.Set("start_address", startAddress)

	return getAddressPoolRecordContext(ctx, d, m)
}

func getAddressPoolRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	startAddr := d.Get("start_address").(string)

	query := map[string]string{
		"startAddress": startAddr,
	}

	response, err := objMgr.GetAddressPool(query)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when getting address pool",
			Detail:   fmt.Sprintf("Error when getting address pool: %v", err),
		})
		return diags
	}

	setIPCAddressPoolResource(d, *response)
	return diags
}

func updateAddressPoolRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	startAddress := d.Get("start_address").(string)
	endAddress := d.Get("end_address").(string)
	addressPoolType := d.Get("type").(string)
	prefixLength := d.Get("prefix_length").(int)
	primaryNetService := d.Get("primary_net_service").(string)
	dhcpOptionSet := d.Get("dhcp_option_set").(string)
	dhcpPolicySet := d.Get("dhcp_policy_set").(string)
	overlapInterfaceIp := d.Get("overlap_interface_ip").(bool)
	name := d.Get("name").(string)
	id := d.Id()

	payload := en.NewAddressPoolPost(en.IPCAddressPoolPost{
		StartAddr:          startAddress,
		EndAddr:            endAddress,
		Name:               name,
		Type:               addressPoolType,
		PrefixLength:       prefixLength,
		PrimaryNetService:  primaryNetService,
		ID:                 id,
		DhcpOptionSet:      dhcpOptionSet,
		OverlapInterfaceIp: overlapInterfaceIp,
		DhcpPolicySet:      dhcpPolicySet,
	})

	_, err = objMgr.UpdateAddressPool(payload)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when updating address pool.",
			Detail:   fmt.Sprintf("Error when updating address pool : %v", err),
		})
		return diags
	}

	return getAddressPoolRecordContext(ctx, d, m)
}

func deleteAddressPoolRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	log.Printf("[DEBUG] %s: Beginning Deletion of address pool", rsSubnetIdString(d))
	container := d.Get("container").(string)
	startAddress := d.Get("start_address").(string)

	_, err := objMgr.DeleteAddressPool(startAddress, container)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete Address pool failed.",
			Detail:   fmt.Sprintf("Delete Address pool  (%s) failed : %s", startAddress, err),
		})
		return diags
	}

	log.Printf("[DEBUG] %s: Deletion of address pool complete", rsSubnetIdString(d))

	return diags
}

func isIPv4(addr string) bool {
	ip := net.ParseIP(addr)
	return ip != nil && ip.To4() != nil
}

func isIPv6(addr string) bool {
	ip := net.ParseIP(addr)
	return ip != nil && ip.To4() == nil
}
