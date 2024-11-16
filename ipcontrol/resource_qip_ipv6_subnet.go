package ipcontrol

import (
	"context"
	"fmt"
	"log"
	"strconv"

	// "strconv"

	"strings"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceQipIPv6Subnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: createQipIPv6SubnetRecord,
		ReadContext:   getQipIPv6SubnetRecord,
		UpdateContext: updateQipIPv6SubnetRecord,
		DeleteContext: deleteQipIPv6SubnetRecord,

		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_prefix_length": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
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
			"block_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"block_prefix_length": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_reverse_zone": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func createQipIPv6SubnetRecord(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] %s: Beginning network block Creation", rsSubnetIdString(d))
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	var err error
	var diags diag.Diagnostics
	subnet := getQipIPv6SubnetFromResourceData(d)

	log.Println("[DEBUG] QIP IPv6 Subnet: " + fmt.Sprintf("%v", subnet))

	err = objMgr.CreateQipIpv6Subnet(subnet)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Creation of QIP IPv6 Subnet failed",
			Detail:   fmt.Sprintf("Creation of Subnet (%s) failed: %s", subnet.SubnetAddress, err),
		})
		return diags
	}

	d.SetId(subnet.SubnetAddress)

	return getQipIPv6SubnetRecord(ctx, d, m)
}

func getQipIPv6SubnetRecord(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	var diags diag.Diagnostics

	orgName := strings.TrimSpace(d.Get("org_name").(string))
	subnetAddress := strings.TrimSpace(d.Get("subnet_address").(string))
	subnetPrefixLength := d.Get("subnet_prefix_length").(int)

	query := map[string]string{
		"orgName":        orgName,
		"addressVersion": "6",
		"subnetAddress":  subnetAddress,
		"prefixLength":   strconv.Itoa(subnetPrefixLength),
	}

	response, err := objMgr.GetQipIPv6Subnet(query)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Getting QIP IPv6 Subnet Failed",
			Detail:   fmt.Sprintf("Getting QIP IPv6 Subnet (%s) failed : %s", subnetAddress, err),
		})
		return diags
	}

	flattenQipIPv6Subnet(d, response)

	log.Printf("[DEBUG] %s: Completed reading subnet block", rsSubnetIdString(d))

	return nil
}

func updateQipIPv6SubnetRecord(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	var err error
	var diags diag.Diagnostics

	orgName := strings.TrimSpace(d.Get("org_name").(string))
	subnetAddress := strings.TrimSpace(d.Get("subnet_address").(string))
	prefixLength := d.Get("subnet_prefix_length").(int)
	subnetName := strings.TrimSpace(d.Get("subnet_name").(string))

	subnet := en.NewQipIPv6SubnetModify(en.QipIPv6SubnetModify{
		OrgName:        orgName,
		AddressVersion: 6,
		SubnetAddress:  subnetAddress,
		PrefixLength:   prefixLength,
		SubnetName:     subnetName,
	})

	_, err = objMgr.UpdateQipIPv6Subnet(subnet)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Updating of QIP IPv6 Subnet failed",
			Detail:   fmt.Sprintf("Updating QIP IPv6 Subnet by Id (%s) failed : %s", d.Id(), err),
		})
		return diags
	}

	return getQipIPv6SubnetRecord(ctx, d, m)

}

func deleteQipIPv6SubnetRecord(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	var diags diag.Diagnostics

	orgName := strings.TrimSpace(d.Get("org_name").(string))
	subnetAddress := strings.TrimSpace(d.Get("subnet_address").(string))
	prefixLength := d.Get("subnet_prefix_length").(int)

	query := map[string]string{
		"orgName":        orgName,
		"subnetAddress":  subnetAddress,
		"prefixLength":   strconv.Itoa(prefixLength),
		"addressVersion": "6",
	}

	log.Println("[DEBUG] Subnet delete: " + fmt.Sprintf("%v", query))

	err := objMgr.DeleteQipIPv6Subnet(query)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Deletion of QIP IPv6 Subnet failed",
			Detail:   fmt.Sprintf("Deleting QIP IPv6 Subnet (%s) failed : %s", subnetAddress, err),
		})
		return diags
	}

	d.SetId(subnetAddress)

	return diags
}

func getQipIPv6SubnetFromResourceData(d *schema.ResourceData) *en.QipIPv6Subnet {

	orgName := strings.TrimSpace(d.Get("org_name").(string))
	subnetAddress := strings.TrimSpace(d.Get("subnet_address").(string))
	blockAddress := strings.TrimSpace(d.Get("block_address").(string))
	prefixLength := d.Get("block_prefix_length").(int)
	subnetPrefixLength := d.Get("subnet_prefix_length").(int)
	createReverseZone := d.Get("create_reverse_zone").(bool)
	poolName := strings.TrimSpace(d.Get("pool_name").(string))
	blockName := strings.TrimSpace(d.Get("block_name").(string))
	subnetName := strings.TrimSpace(d.Get("subnet_name").(string))

	return en.NewQipIPv6Subnet(en.QipIPv6Subnet{
		OrgName:            orgName,
		SubnetAddress:      subnetAddress,
		SubnetName:         subnetName,
		PrefixLength:       prefixLength,
		SubnetPrefixLength: subnetPrefixLength,
		BlockAddress:       blockAddress,
		CreateSubnet:       "SPECIFIC",
		AlgorithmType:      "BEST_FIT_FROM_START",
		AddressVersion:     6,
		CreateReverseZone:  createReverseZone,
		PoolName:           poolName,
		BlockName:          blockName,
	})
}
