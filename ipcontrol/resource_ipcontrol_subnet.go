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

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: createSubnetRecordContext,
		ReadContext:   getSubnetRecordContext,
		UpdateContext: updateSubnetRecordContext,
		DeleteContext: deleteSubnetRecordContext,

		Schema: map[string]*schema.Schema{
			"container": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the container that will hold the block.",
			},
			"rawcontainer": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true to pass the container parameter through to the API without prefixing.",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The address block to allocate.",
			},
			"address_version": {
				ForceNew:    true,
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4,
				Description: "The version of IP Address. Choose 4 for IPV4 or 6 for IPV6. Defaults to 4.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Block Type for the block If not specified, a block type of Any is assumed.",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The size parameter represents the subnet mask or prefix length of the address block in CIDR notation. For IPv4, this is typically a value between 0 and 32 (e.g., 24 for 255.255.255.0). For IPv6, the size value is larger due to the increased address space. IPv6 prefix lengths commonly range between /48 to /128, with /64 often used as the standard size for a single subnet.",
			},
			"dns_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the dns domain that will hold the block.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the block.",
			},
			"block_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The current status of the block. Accepted values are: Deployed, FullyAssigned, Reserved, Aggregate",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify the type of Cloud Provider. Currently one of: AWS, Azure, Cisco ACI, Cisco DNA Center, CloudBolt, OpenStack, ServiceNow, VMware.",
			},
			"cloud_object_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of this object as it is known in the cloud environment.",
			},
		},
	}
}

func createSubnetRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] %s: Beginning network block Creation", rsSubnetIdString(d))
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	address := strings.TrimSpace(d.Get("address").(string))
	container := strings.TrimSpace(d.Get("container").(string))
	rawContainer := d.Get("rawcontainer").(bool)
	size := d.Get("size").(int)
	status := strings.TrimSpace(d.Get("block_status").(string))
	addressVersion := d.Get("address_version").(int)
	blockType := strings.TrimSpace(d.Get("type").(string))
	DNSDomain := strings.TrimSpace(d.Get("dns_domain").(string))
	name := strings.TrimSpace(d.Get("name").(string))
	cloudType := strings.TrimSpace(d.Get("cloud_type").(string))
	cloudObjectId := strings.TrimSpace(d.Get("cloud_object_id").(string))

	var err error
	subnet := en.NewSubnetPost(en.IPCSubnetPost{
		Container:      container,
		RawContainer:   rawContainer,
		Address:        address,
		AddressVersion: addressVersion,
		Type:           blockType,
		Size:           size,
		DNSDomain:      DNSDomain,
		Name:           name,
		BlockStatus:    status,
		CloudType:      cloudType,
		CloudObjectId:  cloudObjectId,
	})

	log.Println("[DEBUG] Subnet post: " + fmt.Sprintf("%v", subnet))

	_, err = objMgr.CreateSubnet(subnet)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create subnet failed",
			Detail:   fmt.Sprintf("Create Subnet block (%s) failed : %s", address, err),
		})
		return diags
	}

	// d.SetId(address)
	log.Printf("[DEBUG] SubnetId: '%s': Creation on network block complete", rsSubnetIdString(d))

	return getSubnetRecordContext(ctx, d, m)
}

func getSubnetRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	log.Printf("[DEBUG] %s: Reading the required subnet block", rsSubnetIdString(d))

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

	flattenIPCSubnet(d, response)
	log.Printf("[DEBUG] %s: Completed reading subnet block", rsSubnetIdString(d))

	return diags
}

func updateSubnetRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	name := d.Get("name").(string)
	size := d.Get("size").(int)

	address := d.Get("address").(string)
	cloudType := d.Get("cloud_type").(string)
	cloudObjectId := d.Get("cloud_object_id").(string)

	_, err = objMgr.UpdateSubnet(address, name, size, cloudType, cloudObjectId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Update Subnet failed",
			Detail:   fmt.Sprintf("Update Subnet block (%s) failed : %s", address, err),
		})
		return diags
	}

	return getSubnetRecordContext(ctx, d, m)

}

func deleteSubnetRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	log.Printf("[DEBUG] %s: Beginning Deletion of network block", rsSubnetIdString(d))
	size := d.Get("size").(int)
	address := d.Get("address").(string)

	_, err := objMgr.DeleteSubnetByIdRef(address, strconv.Itoa(size))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete Subnet failed",
			Detail:   fmt.Sprintf("Delete Subnet block (%s) failed : %s", address, err),
		})
		return diags
	}

	d.SetId(address)
	log.Printf("[DEBUG] %s: Deletion of network block complete", rsSubnetIdString(d))

	return diags
}

type rsSubnetIdStringInterface interface {
	Id() string
}

func rsSubnetIdString(d rsSubnetIdStringInterface) string {
	id := d.Id()
	if id == "" {
		id = "<new resource>"
	}
	return fmt.Sprintf("ID = %s", id)
}
