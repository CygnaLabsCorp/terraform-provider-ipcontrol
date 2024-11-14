package ipcontrol

import (
	"fmt"
	"log"
	"strconv"

	// "strconv"

	"strings"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Create: createSubnetRecord,
		Read:   getSubnetRecord,
		Update: updateSubnetRecord,
		Delete: deleteSubnetRecord,

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

func createSubnetRecord(d *schema.ResourceData, m interface{}) error {

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

	// we demand all the create/reserveIps logic to the CAA
	_, err = objMgr.CreateSubnet(subnet)
	if err != nil {
		return err
	}

	// the Id comes back from the CreateSubnet
	// d.SetId(subnet.Address)

	log.Printf("[DEBUG] SubnetId: '%s': Creation on network block complete", rsSubnetIdString(d))

	// now pull the resource up after deployment via the ID
	// return getSubnetRecord(d, m)
	return getSubnetRecord(d, m)
}

func getSubnetRecord(d *schema.ResourceData, m interface{}) error {
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
		return err
	}

	//d.SetId(strconv.Itoa(response.ID))
	//setIPCSubnetDataData(d, response)

	//d.SetId(obj.Id)

	// setting computed properties from returned object JSON
	// container := strings.Join(obj.Container, ",")
	//cc.FormatContainer(&container)
	// d.Set("container", response.Container)
	// d.Set("address", response.BlockAddr)
	// d.Set("type", response.BlockType)
	// d.Set("size", response.BlockSize)
	// d.Set("dns_domain", response.Subnet.ForwardDomains)
	// d.Set("name", response.BlockName)
	// d.Set("block_status", response.BlockStatus)
	// d.Set("cloud_type", response.CloudType)
	// d.Set("cloud_object_id", response.CloudObjectID)
	flattenIPCSubnet(d, response)

	log.Printf("[DEBUG] %s: Completed reading subnet block", rsSubnetIdString(d))

	return nil
}

func updateSubnetRecord(d *schema.ResourceData, m interface{}) error {
	// Warning or errors can be collected in a slice type
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
		return err
	}

	return getSubnetRecord(d, m)

}

func deleteSubnetRecord(d *schema.ResourceData, m interface{}) error {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	log.Printf("[DEBUG] %s: Beginning Deletion of network block", rsSubnetIdString(d))
	size := d.Get("size").(int)
	address := d.Get("address").(string)

	refRes, err := objMgr.DeleteSubnetByIdRef(address, strconv.Itoa(size))
	if err != nil {
		return err
	}
	// return empty string "" if deleted!
	d.SetId(refRes)
	log.Printf("[DEBUG] %s: Deletion of network block complete", rsSubnetIdString(d))

	return nil
}

type rsSubnetIdStringInterface interface {
	Id() string
}

func rsSubnetIdString(d rsSubnetIdStringInterface) string {
	id := d.Id()
	if id == "" {
		id = "<new resource>"
	}
	return fmt.Sprintf("diamondip_subnet (ID = %s)", id)
}
