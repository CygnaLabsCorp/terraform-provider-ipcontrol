package ipcontrol

import (
	"context"
	"fmt"
	"log"

	// "regexp"
	"strconv"
	"strings"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	objMgr "terraform-provider-ipcontrol/ipcontrol/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetsRead,
		Schema: map[string]*schema.Schema{
			/* cloud type is mandatory */
			"cloudtype": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
				ValidateFunc: validation.StringInSlice([]string{
					string("azure"),
					string("ipam"),
				}, true),
				Description: "The IPControl Cloud Type.",
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenantid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			/* you may just filter by parent container or path if just IPAM
			   if 'container' is specified then all the cloud's eventual containers' modelling
			   parameters specified will be ignored (i.e. rg/vnet ...)
			   However if the cloudtype == 'ipam' then container must be specified */
			"container": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"blocktype": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IPControl Block Type of the subnet, if not specified will default to 'Any'.",
				Default:     "Any",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"blockstatus": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Deployed",
				Description: "The IPControl Block Status of the subnet, if not specified will default to 'Deployed'.",
				ValidateFunc: validation.StringInSlice([]string{
					string("deployed"),
					string("reserved"),
					string("aggregate"),
					string("free"),
				}, true),
			},
			"resourcegroup": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtualnetwork": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"subnets": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenantid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"container": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						/* if cloudtype is null, then it'll just operate with normal IPAM topology without any logical mapping of CSP terms */
						"cloudtype": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"blocktype": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"blockstatus": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						// optional/conditional data arguments
						"cloudobjid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						/* Azure Selector Parameters */
						"resourcegroup": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtualnetwork": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"location": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						// computed data arguments
						"tenantname": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastupdated": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		}, // schema

	}
}

func dataSourceSubnetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	connector := m.(*objMgr.Connector)
	objMgr := objMgr.NewObjectManager(connector)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	// collect input properties (from .tf file data element) that will be used to get the object
	cidr := strings.TrimSpace(d.Get("cidr").(string))
	tenantId := strings.TrimSpace(d.Get("tenantid").(string))
	cloudType := strings.ToLower(d.Get("cloudtype").(string))
	name := strings.TrimSpace(d.Get("name").(string))
	blocktype := strings.TrimSpace(strings.ToLower(d.Get("blocktype").(string)))
	blockstatus := strings.ToLower(d.Get("blockstatus").(string))
	container := strings.TrimSpace(d.Get("container").(string))

	// create parames slice
	parMap := make(map[string]string)

	/* Azure Selection Logic */
	if cloudType == "azure" {
		/* the following selector parameters are then expected
		 * resourceGroup or Location
		 * VNET
		 */

		rg := strings.TrimSpace(d.Get("resourcegroup").(string))
		lc := strings.TrimSpace(d.Get("location").(string))
		vn := strings.TrimSpace(d.Get("virtualnetwork").(string))

		// Build proper Selector element
		// rg is required unless location is specified, in which case the latter will take precedence
		if lc != "" {
			parMap["location"] = lc
		} else if rg != "" {
			parMap["resourcegroup"] = rg
		} else {

			// ERROR - none of them were specified
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "dataSourceSubnetsRead - [Azure] Exporting Subnets failed missing selector arguments",
				Detail:   fmt.Sprintf("One of the selector arguments 'location' or 'resourcegroup' is required, but no definition of any of them was found."),
			})

			return diags
		}

		if vn != "" {
			parMap["virtualnetwork"] = vn
		} else {

			// ERROR - vnet was not specified
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "dataSourceSubnetsRead - [Azure] Exporting Subnets failed missing selector arguments",
				Detail:   fmt.Sprintf("The selector arguments 'virtualnetwork' is required, but no definition was found."),
			})

			return diags
		}

	}

	parMap["tenantid"] = tenantId
	parMap["cloudtype"] = cloudType

	if cidr != "" {
		parMap["cidr"] = cidr
	}
	if container != "" {
		parMap["container"] = container
	}
	if blocktype != "" {
		parMap["blocktype"] = blocktype
	}
	if blockstatus != "" {
		parMap["blockstatus"] = blockstatus
	}
	if name != "" {
		parMap["name"] = name
	}

	params := en.Params{}
	for k, v := range parMap {
		params[k] = v
	}

	var response *[]en.Subnet
	response, err = objMgr.ExportSubnets(params)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Getting Subnet failed",
			Detail:   fmt.Sprintf("Getting Subnet block (%s) failed : %s", cidr, err),
		})
		return diags
	}

	/* log examples:

	   log.Println("[DEBUG] Subnet Objects: " + fmt.Sprintf("%s", response))
	   for i, subnet := range *response {
	       log.Println("[DEBUG] ONe Test: " + fmt.Sprintf("%v %s", i, string(subnet.Cidr)))
	   }

	*/

	if response == nil { // || (reflect.TypeOf(response) == reflect.TypeOf(data) && len(response.([](cc.Subnet))) == 0) {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "API returns a nil/empty id",
			Detail:   fmt.Sprintf("API returns a nil/empty subnet response. Getting Subnet block (%s) failed", cidr),
		})

		return diags

	}

	// *** always set the ID ***
	// d.SetId("xxxx")
	// we've to definitievely store an ID on the returned resouceData object pointer, and in such operations we just  may use a timestamp,
	// we do not need a contextual ID like on CRUD operations given the nature of the "export" function ...
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	subnets := flattenSubnetsData(response)
	log.Println("[DEBUG] Subnets: " + fmt.Sprintf("%v", subnets))

	if err := d.Set("subnets", subnets); err != nil {
		return diag.FromErr(err)
	}

	// log.Println("[DEBUG] Subnet Object: " + fmt.Sprintf("%v", response))

	return diags
}

/* flattenSubnetsData
 * returns an array of objects for straight JSON marshalling from the resourceData object returned
 */
func flattenSubnetsData(subnets *[]en.Subnet) []interface{} {
	if subnets != nil {
		sbs := make([]interface{}, len(*subnets), len(*subnets))

		for i, subnet := range *subnets {
			sb := make(map[string]interface{})

			sb["id"] = string(subnet.Id)
			sb["cidr"] = string(subnet.Cidr)
			sb["name"] = string(subnet.Name)
			sb["gateway"] = string(subnet.Gateway)
			sb["blocktype"] = string(subnet.Blocktype)
			sb["blockstatus"] = string(subnet.Blockstatus)
			sb["description"] = string(subnet.Description)
			sb["cloudtype"] = string(subnet.Cloudtype)
			sb["cloudobjid"] = string(subnet.Cloudobjid)
			sb["tenantid"] = string(subnet.Tenantid)
			sb["tenantname"] = string(subnet.Tenantname)
			sb["lastupdated"] = string(subnet.Lastupdated)
			sb["container"] = string(subnet.Container)

			sbs[i] = sb
		}

		return sbs
	}

	return make([]interface{}, 0)
}
