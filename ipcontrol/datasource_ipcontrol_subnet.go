package ipcontrol

import (

	// "regexp"

	"strings"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	objMgr "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSubnetsRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
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
		}, // schema

	}
}

func dataSourceSubnetsRead(d *schema.ResourceData, m interface{}) error {
	connector := m.(*objMgr.Connector)
	objMgr := objMgr.NewObjectManager(connector)

	// Warning or errors can be collected in a slice type
	//var diags diag.Diagnostics
	var err error

	username := strings.TrimSpace(d.Get("username").(string))
	password := strings.TrimSpace(d.Get("password").(string))
	//container := strings.TrimSpace(d.Get("container").(string))
	address := strings.TrimSpace(d.Get("address").(string))
	typeSubnet := strings.TrimSpace(d.Get("type").(string))
	size := strings.TrimSpace(d.Get("size").(string))

	// create parames slice
	parMap := make(map[string]string)

	parMap["username"] = username
	parMap["password"] = password
	//parMap["container"] = container
	parMap["address"] = address
	parMap["typeSubnet"] = typeSubnet
	parMap["size"] = size

	params := en.Params{}
	for k, v := range parMap {
		params[k] = v
	}

	var response *en.IPCSubnet
	response, err = objMgr.GetSubnetByIdRef("1")
	if err != nil {
		// diags = append(diags, diag.Diagnostic{
		// 	Severity: diag.Error,
		// 	Summary:  "Getting Subnet failed",
		// 	Detail:   fmt.Sprintf("Getting Subnet block (%s) failed : %s", address, err),
		// })
		return err
	}

	if response == nil { // || (reflect.TypeOf(response) == reflect.TypeOf(data) && len(response.([](cc.Subnet))) == 0) {

		// diags = append(diags, diag.Diagnostic{
		// 	Severity: diag.Error,
		// 	Summary:  "API returns a nil/empty id",
		// 	Detail:   fmt.Sprintf("API returns a nil/empty subnet response. Getting Subnet block (%s) failed", address),
		// })

		return nil

	}

	// *** always set the ID ***
	// d.SetId("xxxx")
	// we've to definitievely store an ID on the returned resouceData object pointer, and in such operations we just  may use a timestamp,
	// we do not need a contextual ID like on CRUD operations given the nature of the "export" function ...
	// d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	// subnets := flattenSubnetsData(response)
	// log.Println("[DEBUG] Subnets: " + fmt.Sprintf("%v", subnets))

	// if err := d.Set("subnets", subnets); err != nil {
	// 	return diag.FromErr(err)
	// }

	// log.Println("[DEBUG] Subnet Object: " + fmt.Sprintf("%v", response))

	return nil
}
