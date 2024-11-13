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
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_version": {
				ForceNew: true,
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
			},
		},
	}
}

func createSubnetRecord(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] %s: Beginning network block Creation", rsSubnetIdString(d))
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	// trimmed strings
	container := strings.TrimSpace(d.Get("container").(string))
	cc.FormatContainer(&container)
	address := strings.TrimSpace(d.Get("address").(string))
	typeSubnet := strings.TrimSpace(d.Get("type").(string))
	size := d.Get("size").(int)
	name := d.Get("name").(string)
	version := d.Get("address_version").(int)

	log.Printf("[DEBUG] SubnetId: '%s': Creation on network block complete", rsSubnetIdString(d))

	var subnet *en.IPCSubnetPost
	var err error

	// we demand all the create/reserveIps logic to the CAA
	subnet, err = objMgr.CreateSubnet(container, address, typeSubnet, size, name, version)
	if err != nil {
		return err
	}

	// the Id comes back from the CreateSubnet
	d.SetId(subnet.Address)

	log.Printf("[DEBUG] SubnetId: '%s': Creation on network block complete", rsSubnetIdString(d))

	// now pull the resource up after deployment via the ID
	// return getSubnetRecord(d, m)
	return nil
}

func getSubnetRecord(d *schema.ResourceData, m interface{}) error {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	log.Printf("[DEBUG] %s: Reading the required subnet block", rsSubnetIdString(d))

	obj, err := objMgr.GetSubnetByIdRef(string(d.Id()))
	if err != nil {
		return err
	}
	//d.SetId(obj.Id)

	// setting computed properties from returned object JSON
	// d.Set("username", string(obj.Username))
	// d.Set("password", string(obj.Password))
	container := strings.Join(obj.Container, ",")
	cc.FormatContainer(&container)
	d.Set("container", container)
	d.Set("address", string(obj.Address))
	d.Set("size", obj.Size)
	d.Set("name", string(obj.Name))

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

	idRef := d.Id()

	_, err = objMgr.UpdateSubnet(idRef, name, size)

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

	refRes, err := objMgr.DeleteSubnetByIdRef(string(d.Id()), strconv.Itoa(size))
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
