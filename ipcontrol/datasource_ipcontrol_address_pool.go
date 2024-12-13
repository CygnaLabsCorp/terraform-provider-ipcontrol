package ipcontrol

import (
	"context"
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAddressPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAddressPoolRead,

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
				Computed:    true,
				Description: "The end address of the address pool.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of address pool.",
			},
			"prefix_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size parameter represents the subnet mask or prefix length of the address block in CIDR notation. Required for IPv6, the size value is larger due to the increased address space. IPv6 prefix lengths commonly range between /48 to /128, with /64 often used as the standard size for a single subnet.",
			},
			"primary_net_service": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the DHCP server that will serve addresses from this pool.",
			},
			"overlap_interface_ip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to allow a DHCPv6 pool to overlap an interface address.",
			},
			"dhcp_policy_set": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of a Policy Set used with this pool.",
			},
			"dhcp_option_set": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of a Option Set used with this pool.",
			},
			"name": {
				Type:        schema.TypeString,
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
	}
}

func dataSourceAddressPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	startAddr := d.Get("start_address").(string)
	container := d.Get("container").(string)

	query := map[string]string{
		"startAddress": startAddr,
	}

	if container != "" {
		query["container"] = container
	}

	response, err := objMgr.GetAddressPool(query)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when reading address pool",
			Detail:   fmt.Sprintf("Error when reading address pool: %v", err),
		})
		return diags
	}

	setIPCAddressPoolResource(d, *response)
	return diags
}
