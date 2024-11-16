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

func dataSourceQipIPv6Subnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQipIPv6SubnetRead,
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"block_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_prefix_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceQipIPv6SubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	orgName := strings.TrimSpace(d.Get("org_name").(string))
	subnetAddress := strings.TrimSpace(d.Get("subnet_address").(string))
	subnetPrefixLength := d.Get("subnet_prefix_length").(int)
	subnetName := strings.TrimSpace(d.Get("subnet_name").(string))

	query := map[string]string{
		"orgName":        orgName,
		"addressVersion": "6",
	}

	if subnetName != "" {
		query["subnetName"] = subnetName
	} else if subnetAddress != "" && subnetPrefixLength > 0 {
		query["subnetAddress"] = subnetAddress
		query["prefixLength"] = strconv.Itoa(subnetPrefixLength)
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Getting QIP IPv6 Subnet Failed",
			Detail:   "Missing subnet_name and subnet_address/subnet_prefix_length field",
		})
		return diags
	}

	response, err := objMgr.GetQipIPv6Subnet(query)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Getting QIP IPv6 Subnet Failed",
			Detail:   fmt.Sprintf("Getting QIP IPv6 Subnet (%s) failed : %s", response.SubnetAddress, err),
		})
		return diags
	}

	if response == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "API returns a nil/empty QIP IPv6 Subnet",
			Detail:   "API returns a nil/empty subnet response. Getting QIP IPv6 Subnet failed",
		})
		return diags
	}

	flattenQipIPv6Subnet(d, response)
	log.Println("[DEBUG] QIP IPv6 Subnet Object: " + fmt.Sprintf("%v", response))

	return diags
}

func flattenQipIPv6Subnet(d *schema.ResourceData, subnet *en.QipIPv6SubnetGet) {

	d.SetId(subnet.SubnetAddress)
	d.Set("subnet_address", subnet.SubnetAddress)
	d.Set("pool_name", subnet.PoolName)
	d.Set("block_name", subnet.BlockName)
	d.Set("subnet_prefix_length", subnet.SubnetPrefixLength)
	d.Set("subnet_name", subnet.SubnetName)
}
