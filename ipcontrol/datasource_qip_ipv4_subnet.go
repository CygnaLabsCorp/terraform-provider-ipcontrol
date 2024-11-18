package ipcontrol

import (

	// "regexp"

	"context"
	"fmt"
	"log"
	"strings"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceQipIPv4Subnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQipIPv4SubnetRead,
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization Name.",
			},
			"subnet_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IPv4 address of subnet.",
			},
			"subnet_mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Subnet mask of subnet.",
			},
			"network_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv4 address of network.",
			},
			"warning_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Type of warning when the defined percentage of addresses reached.",
			},
			"warning_percent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Percentage of managed addresses. Warning if this percentage is reached. The value range is 0 - 100.",
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of subnet.",
			},
		},
	}
}

func dataSourceQipIPv4SubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	orgName := strings.TrimSpace(d.Get("org_name").(string))
	subnetAddress := strings.TrimSpace(d.Get("subnet_address").(string))

	query := map[string]string{
		"orgName":        orgName,
		"subnetAddress":  subnetAddress,
		"addressVersion": "4",
	}

	response, err := objMgr.GetQipIPv4Subnet(query)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Getting QIP IPv4 Subnet Failed",
			Detail:   fmt.Sprintf("Getting QIP IPv4 Subnet (%s) failed : %s", subnetAddress, err),
		})
		return diags
	}

	if response == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "API returns a nil/empty QIP IPv4 Subnet",
			Detail:   fmt.Sprintf("API returns a nil/empty subnet response. Getting QIP IPv4 Subnet (%s) failed", subnetAddress),
		})
		return diags
	}

	flattenQipIPv4Subnet(d, response)
	log.Println("[DEBUG] Subnet Object: " + fmt.Sprintf("%v", response))

	return nil
}

func flattenQipIPv4Subnet(d *schema.ResourceData, subnet *en.QipIPv4Subnet) {

	d.SetId(subnet.SubnetAddress)
	d.Set("subnet_address", subnet.SubnetAddress)
	d.Set("subnet_mask", subnet.SubnetMask)
	d.Set("network_address", subnet.NetworkAddress)
	d.Set("warning_type", subnet.WarningType)
	d.Set("warning_percent", subnet.WarningPercentage)
	d.Set("subnet_name", subnet.SubnetName)
}
