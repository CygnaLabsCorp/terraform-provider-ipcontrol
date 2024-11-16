package ipcontrol

import (
	"context"
	"fmt"
	"log"

	// "strconv"

	"strings"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceQipIPv4Subnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: createQipIPv4SubnetRecord,
		ReadContext:   getQipIPv4SubnetRecord,
		UpdateContext: updateQipIPv4SubnetRecord,
		DeleteContext: deleteQipIPv4SubnetRecord,

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
			"subnet_mask": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"network_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"warning_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"warning_percent": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func createQipIPv4SubnetRecord(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] %s: Beginning network block Creation", rsSubnetIdString(d))
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	var err error
	var diags diag.Diagnostics
	subnet := getQipIPv4SubnetFromResourceData(d)

	log.Println("[DEBUG] Subnet post: " + fmt.Sprintf("%v", subnet))

	_, err = objMgr.CreateQipIPv4Subnet(subnet)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Creation of QIP IPv4 Subnet failed",
			Detail:   fmt.Sprintf("Creation of Subnet (%s) failed: %s", subnet.SubnetAddress, err),
		})
		return diags
	}

	d.SetId(subnet.SubnetAddress)
	log.Printf("[DEBUG] SubnetId: '%s': Creation on network block complete", rsSubnetIdString(d))

	return getQipIPv4SubnetRecord(ctx, d, m)
}

func getQipIPv4SubnetRecord(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	var diags diag.Diagnostics

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

	flattenQipIPv4Subnet(d, response)

	log.Printf("[DEBUG] %s: Completed reading subnet block", rsSubnetIdString(d))

	return nil
}

func updateQipIPv4SubnetRecord(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	var err error
	var diags diag.Diagnostics
	subnet := getQipIPv4SubnetFromResourceData(d)

	_, err = objMgr.UpdateQipIPv4Subnet(subnet)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Updating of QIP IPv4 Subnet failed",
			Detail:   fmt.Sprintf("Updating QIP IPv4 Subnet by Id (%s) failed : %s", d.Id(), err),
		})
		return diags
	}

	return getQipIPv4SubnetRecord(ctx, d, m)

}

func deleteQipIPv4SubnetRecord(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	var diags diag.Diagnostics

	orgName := strings.TrimSpace(d.Get("org_name").(string))
	subnetAddress := strings.TrimSpace(d.Get("subnet_address").(string))

	query := map[string]string{
		"orgName":        orgName,
		"subnetAddress":  subnetAddress,
		"addressVersion": "4",
	}

	log.Println("[DEBUG] Subnet post: " + fmt.Sprintf("%v", query))

	err := objMgr.DeleteQipIPv4Subnet(query)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Deletion of QIP IPv4 Subnet failed",
			Detail:   fmt.Sprintf("Deleting QIP IPv4 Subnet block by Id (%s) failed : %s", d.Id(), err),
		})
		return diags
	}

	d.SetId(subnetAddress)
	//log.Printf("[DEBUG] %s: Deletion of network block complete", rsSubnetIdString(d))

	return diags
}

func getQipIPv4SubnetFromResourceData(d *schema.ResourceData) *en.QipIPv4Subnet {

	orgName := strings.TrimSpace(d.Get("org_name").(string))
	subnetAddress := strings.TrimSpace(d.Get("subnet_address").(string))
	subnetMask := strings.TrimSpace(d.Get("subnet_mask").(string))
	networkAddress := strings.TrimSpace(d.Get("network_address").(string))
	subnetName := strings.TrimSpace(d.Get("subnet_name").(string))
	warningType := d.Get("warning_type").(int)
	warningPercent := d.Get("warning_percent").(int)

	return en.NewQipIPv4Subnet(en.QipIPv4Subnet{
		OrgName:           orgName,
		SubnetAddress:     subnetAddress,
		SubnetMask:        subnetMask,
		NetworkAddress:    networkAddress,
		SubnetName:        subnetName,
		WarningType:       warningType,
		WarningPercentage: warningPercent,
		AddressVersion:    4,
	})
}
