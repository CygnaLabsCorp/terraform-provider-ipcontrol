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
			},
			"rawcontainer": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_version": {
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

	flattenIPCSubnet(d, response)

	log.Println("[DEBUG] Subnet Object: " + fmt.Sprintf("%v", response))

	return nil
}

func flattenIPCSubnet(d *schema.ResourceData, subnet *en.IPCSubnet) {

	d.SetId(subnet.BlockAddr)
	d.Set("container", subnet.Container[0])
	d.Set("address", subnet.BlockAddr)
	d.Set("type", subnet.BlockType)
	d.Set("size", subnet.BlockSize)
	d.Set("name", subnet.BlockName)
	d.Set("block_status", subnet.BlockStatus)
	d.Set("cloud_type", subnet.CloudType)
	d.Set("cloud_object_id", subnet.CloudObjectID)
}
