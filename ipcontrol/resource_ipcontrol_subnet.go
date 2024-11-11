package ipcontrol

import (
	"fmt"
	"log"

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
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"container": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createSubnetRecord(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] %s: Beginning network block Creation", rsSubnetIdString(d))
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	// trimmed strings
	username := strings.TrimSpace(d.Get("username").(string))
	password := strings.TrimSpace(d.Get("password").(string))
	container := strings.TrimSpace(d.Get("container").(string))
	address := strings.TrimSpace(d.Get("address").(string))
	typeSubnet := strings.TrimSpace(d.Get("type").(string))
	size := strings.TrimSpace(d.Get("size").(string))

	log.Printf("[DEBUG] SubnetId: '%s': Creation on network block complete", rsSubnetIdString(d))

	var subnet *en.IPC_Subnet_Post
	var err error

	// we demand all the create/reserveIps logic to the CAA
	subnet, err = objMgr.CreateSubnet(username, password, container, address, typeSubnet, size)
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
	d.Set("container", strings.Join(obj.Container, ","))
	d.Set("address", string(obj.Address))
	d.Set("type", string(obj.Type))
	d.Set("size", string(obj.Size))

	log.Printf("[DEBUG] %s: Completed reading subnet block", rsSubnetIdString(d))

	return nil
}

func updateSubnetRecord(d *schema.ResourceData, m interface{}) error {
	// Warning or errors can be collected in a slice type
	var err error
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	username := strings.TrimSpace(d.Get("username").(string))
	password := strings.TrimSpace(d.Get("password").(string))
	container := strings.TrimSpace(d.Get("container").(string))
	address := strings.TrimSpace(d.Get("address").(string))
	typeSubnet := strings.TrimSpace(d.Get("type").(string))
	size := strings.TrimSpace(d.Get("size").(string))

	// create params slice
	parMap := make(map[string]string)

	parMap["username"] = username
	parMap["password"] = password
	parMap["container"] = container
	parMap["address"] = address
	parMap["type"] = typeSubnet
	parMap["size"] = size

	params := en.Params{}
	for k, v := range parMap {
		params[k] = v
	}

	idRef := d.Id()

	_, err = objMgr.UpdateSubnet(idRef, params)

	if err != nil {
		return err
	}

	return getSubnetRecord(d, m)

}

func deleteSubnetRecord(d *schema.ResourceData, m interface{}) error {
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	log.Printf("[DEBUG] %s: Beginning Deletion of network block", rsSubnetIdString(d))

	refRes, err := objMgr.DeleteSubnetByIdRef(string(d.Id()))
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
