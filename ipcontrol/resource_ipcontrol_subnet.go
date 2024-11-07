package ipcontrol

import (
	"context"
	"fmt"
	"log"

	// "strconv"
	"errors"
	"regexp"
	"strings"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetCreate,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetUpdate,
		DeleteContext: resourceSubnetDelete,

		Schema: map[string]*schema.Schema{
			// required data arguments
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // if you change this, it'll force re-creation
				Description: "Set parameter value = 'next/(size|h:<nrOfHostsEff>)' to allocate next available cidr of specific size (i.e. next class-C subnet: 'next/24') or next subnet of best size based on required number of effective hosts (i.e. next subnet that can support up to i.e. 40 effecitve hosts' IP addresses: 'next/h:40'). Alternatively,  specific subnet address can also be requested replacing 'next' with desired block address, i.e.  '10.0.0.0/24' or '10.0.0.0/h:40'",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^(\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\/([8-9]|[1-2][0-9]|3[0-2]|h:\d+)))|(next\/([8-9]|[1-2][0-9]|3[0-2]|h:\d+))$`),
					"'cidr' must be a valid block prefix without spaces i.e. '10.0.0.0/24' or, for auto-allocation, provided in the form: 'next/<size>', i.e. 'next/24'. Aleternatively, block size can be also requested by number of effective IP hosts in the form of 'h:<nrOfHostsEff>', i.e. 10.0.0.0/h:40 or next/h:40",
				),
			},
			// mandatory if cloudtype !== ipam
			"tenantid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"container": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloudtype": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
				ValidateFunc: validation.StringInSlice([]string{
					string("azure"),
					string("ipam"),
				}, true),
				Description: "The IPControl Cloud Type.",
			},
			// optional/conditional data arguments
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of your subnet block, defaults to subnet cidr.",
				// retain existing value if updating with empty name
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" && new == "" {
						return true
					}
					return false
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of your subnet block.",
			},
			"reserveips": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     0,
				Description: "The IP's you want to reserve.",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^from-(start|end)\/[1-9][0-9]*$`),
					"Reserve IPs (reserveips) parameter is invalid. Accepted only at creation phase. Expected format: 'from-[start|end]/<nrOfIps>', where <nrOfIps> must be > 0. Example: 'from-start/4'.",
				),
			},
			// parent block, if specified after creation will force a re-deployment
			"parentblock": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The parent address block in cidr format to allocate from.",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^(\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\/([8-9]|[1-2][0-9]|3[0-2])))$`),
					"'parentblock' must be single valid block cide without spaces i.e. '10.0.0.0/16'.",
				),
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
				}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"gateway": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "gateway ip address of your network block.By default first IPv4 address is set as gateway address.",
				/* accepted as IP address (on create or update) or next/from-[start|end]-offset/number (only on create) */
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|(.+\/from-(start|end)-offset\/[1-9][0-9]*)$`),
					"Gateway must be valid IP Address, i.e. '10.0.0.1' or - format available only on Create - in the form: 'next/from-[start|end]-offset/<offset>', where <offset> must be > 0. Example: 'next/from-start-offset/1'.",
				),
			},
			/* Azure Selector Parameters */
			"resourcegroup": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"virtualnetwork": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			// computed data arguments
			"tenantname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloudobjid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"lastupdated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		/* cross-check input paramters depending on the specified cloud types */
		CustomizeDiff: customdiff.Sequence(

			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				if v, ok := diff.GetOk("cloudtype"); ok {
					/* ** cloud type --> 'ipam':
					   + container parameter is required
					   + parentblock is not supported

					   ** cloud_type --> 'azure:
					   + resourcegroup or location must be specified
					   + virtualnetwork must be specified
					   + container is not accepted and should not be specified

					*/

					if strings.ToLower(v.(string)) == "ipam" {
						// container parameter is required
						if v, ok := diff.GetOk("container"); ok && v.(string) != "" {
							return nil
						} else {
							return errors.New(`Cloud Type "IPAM" requires non-empty 'container' input parameter specified.`)
						}
						// parentblock parameter must not be specified
						if v, ok := diff.GetOk("parentblock"); !ok && v.(string) != "" {
							return errors.New(`Cloud Type "IPAM" cannot accept 'parentblock' input parameter.`)
						}
					} else {
						// this is cloud type, then we throw error in case of container value specified, though the filter below
						// does not capture ""
						// then the create method will filter the empty value from the msg params
						if v, ok := diff.GetOk("container"); ok && v.(string) != "" {
							ct := strings.TrimSpace(diff.Get("cloudtype").(string))
							errmsg := fmt.Sprintf(`Parameter 'container' is only accepted for cloud type = 'IPAM' (Found cloudtype='%s')`, ct)
							return errors.New(errmsg)
						}

						// validate Azure input params rGroup/Location and vNet
						if strings.ToLower(v.(string)) == "azure" {

							rg := strings.TrimSpace(diff.Get("resourcegroup").(string))
							lc := strings.TrimSpace(diff.Get("location").(string))
							vn := strings.TrimSpace(diff.Get("virtualnetwork").(string))
							if rg == "" && lc == "" {
								return errors.New(`'resourcegroup' or 'location' parameter must be specified with 'Azure' cloud type.`)
							}
							if vn == "" {
								return errors.New(`Parameter 'virtualnetwork' is required with 'Azure' cloud type.`)
							}

						}
					}
				}

				return nil

			},
		),
	}
}

func resourceSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	log.Printf("[DEBUG] %s: Beginning network block Creation", rsSubnetIdString(d))
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	/* On Create, following parameters are supported:
	 *
	 * cidr (mandatory)
	 * cloud-related container selector (tenantid, resGroup, vNet for Azure)
	 * blocktype (optional, defaults to 'Any' if not specified)
	 * blockstatus <Reserved/Deployed> (optinal, defaults to 'Deployed' if not specified)
	 * name (optional, defaults to cidr)
	 * description (optional)
	 * gateway <IpAddress|next/from-[start|end]-offset> (optional)
	 * tenantid (mandatory)
	 * cloudtype (mandatory)
	 * reverve_ips
	 */

	// trimmed strings
	cidr := strings.TrimSpace(d.Get("cidr").(string))
	blocktype := strings.TrimSpace(strings.ToLower(d.Get("blocktype").(string)))
	name := strings.TrimSpace(d.Get("name").(string))
	description := strings.TrimSpace(d.Get("description").(string))
	tenantId := strings.TrimSpace(d.Get("tenantid").(string))

	// pre-validated with 'ValidateFunc'
	blockstatus := strings.ToLower(d.Get("blockstatus").(string))
	gateway := d.Get("gateway").(string)
	cloudType := strings.ToLower(d.Get("cloudtype").(string))
	container := strings.TrimSpace(d.Get("container").(string))
	parentblock := d.Get("parentblock").(string)
	reserveips := d.Get("reserveips").(string)

	// create parames slice
	parMap := make(map[string]string)

	/* Azure Selection Logic */
	if cloudType == "azure" {
		/* the following selector parameters are then expected
		 * resourceGroup or Region
		 * VNET
		 */

		rg := strings.TrimSpace(d.Get("resourcegroup").(string))
		lc := strings.TrimSpace(d.Get("location").(string))
		vn := strings.TrimSpace(d.Get("virtualnetwork").(string))

		// Build proper l1Selector element
		// rg is required unless location is specified, in which ase RG will be discared
		if lc != "" {
			parMap["location"] = lc
		} else if rg != "" {
			parMap["resourcegroup"] = rg
		} else {

			// ERROR - none of them were defined
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "resourceSubnetCreate - [Azure] Creating Subnet failed missing selector arguments",
				Detail:   fmt.Sprintf("One of the selector arguments 'location' or 'resourcegroup' is required, but no definition was found."),
			})

			return diags
		}

		if vn != "" {
			parMap["virtualnetwork"] = vn
		} else {

			// ERROR - none of them were defined
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "resourceSubnetCreate - [Azure] Creating Subnet failed missing selector arguments",
				Detail:   fmt.Sprintf("The selector arguments 'virtualnetwork' is required, but no definition was found."),
			})

			return diags
		}

	}

	parMap["tenantid"] = tenantId
	parMap["cloudtype"] = cloudType

	if blocktype != "" {
		parMap["blocktype"] = blocktype
	}
	if blockstatus != "" {
		parMap["blockstatus"] = blockstatus
	}
	if description != "" {
		parMap["description"] = description
	}
	if parentblock != "" {
		parMap["parentblock"] = parentblock
	}
	if reserveips != "" {
		parMap["reserveips"] = reserveips
	}
	if container != "" {
		parMap["container"] = container
	}

	params := en.Params{}
	for k, v := range parMap {
		params[k] = v
	}

	var subnet *en.Subnet
	var err error

	// we demand all the create/reserveIps logic to the CAA
	subnet, err = objMgr.CreateSubnet(cidr, name, gateway, params)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "resourceSubnetCreate - Creation of Subnet failed",
			Detail:   fmt.Sprintf("resourceSubnetCreate - Creation of Subnet failed: %s", err),
		})
		return diags
	}

	// the Id comes back from the CreateSubnet
	d.SetId(subnet.Id)

	log.Printf("[DEBUG] SubnetId: '%s': Creation on network block complete", rsSubnetIdString(d))

	// now pull the resource up after deployment via the ID
	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	log.Printf("[DEBUG] %s: Reading the required subnet block", rsSubnetIdString(d))

	obj, err := objMgr.GetSubnetByIdRef(string(d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "resourceSubnetRead - Getting of Subnet failed",
			Detail:   fmt.Sprintf("Getting Subnet block by Id (%s) failed : %s", d.Id(), err),
		})
		return diags
	}
	d.SetId(obj.Id)

	// setting computed properties from returned object JSON
	d.Set("name", string(obj.Name))
	// d.Set("cloudobjid", obj.CloudObjId)
	d.Set("tenantname", string(obj.Tenantname))
	d.Set("cidr", string(obj.Cidr))
	// d.Set("gateway", obj.Gateway) // i.e. Gateway     string `json:"gateway,omitempty"`

	/* we need to pull these attributes back from the response given they may just be default hence not defined in the create phase
	 * then Read will report what are the deployed effective parameter values */
	d.Set("blocktype", string(obj.Blocktype))
	d.Set("blockstatus", string(obj.Blockstatus))
	d.Set("description", string(obj.Description))
	d.Set("gateway", string(obj.Gateway))
	d.Set("cloudobjid", string(obj.Cloudobjid))
	// retain location from the resource input element if it was passed
	d.Set("lastupdated", string(obj.Lastupdated))
	d.Set("container", string(obj.Container))

	log.Printf("[DEBUG] %s: Completed reading subnet block", rsSubnetIdString(d))

	return diags
}

func resourceSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	blockstatus := d.Get("blockstatus").(string)

	tenantid := strings.TrimSpace(d.Get("tenantid").(string))
	name := strings.TrimSpace(d.Get("name").(string))
	description := strings.TrimSpace(d.Get("description").(string))
	cloudobjid := strings.TrimSpace(d.Get("cloudobjid").(string))
	cidr := strings.TrimSpace(d.Get("cidr").(string))

	// *** Note: gateway here can only accept the IP address form here on PUT
	gateway := d.Get("gateway").(string)
	invalidGw := regexp.MustCompile(`^.+\/from-(start|end)-offset\/[1-9][0-9]$`)

	if invalidGw.MatchString(gateway) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("resourceSubnetUpdate - Updating of Subnet %s failed.", cidr),
			Detail:   fmt.Sprintf("Invalid Gateway input param value on UPDATE: %s", gateway),
		})
		return diags
	}

	/* the following parameters only are accepted for update :
	   blockstatus := d.Get("blockstatus").(string) [only toggle between reserved / deployed]
	   name := d.Get("name").(string)
	   description := d.Get("description").(string)
	   gateway := d.Get("gateway").(string)
	   cloudObjId cloudobjid

	   forceRenew:
	   cidr
	   tenantid
	   parentblock
	   cloudtype
	*/

	// create params slice
	parMap := make(map[string]string)

	if blockstatus != "" {
		parMap["blockstatus"] = blockstatus
	}

	if name != "" {
		parMap["name"] = name
	}

	if description != "" {
		parMap["description"] = description
	}

	if gateway != "" {
		parMap["gateway"] = gateway
	}

	if cloudobjid != "" {
		parMap["cloudobjid"] = cloudobjid
	}

	if tenantid != "" {
		parMap["tenantid"] = tenantid
	}

	params := en.Params{}
	for k, v := range parMap {
		params[k] = v
	}

	idRef := d.Id()

	_, err = objMgr.UpdateSubnet(idRef, params)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "resourceSubnetUpdate - Updating of Subnet failed",
			Detail:   fmt.Sprintf("Updating Subnet block by Id (%s) failed : %s", d.Id(), err),
		})
		return diags
	}

	return resourceSubnetRead(ctx, d, m)

}

func resourceSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	log.Printf("[DEBUG] %s: Beginning Deletion of network block", rsSubnetIdString(d))

	refRes, err := objMgr.DeleteSubnetByIdRef(string(d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "resourceSubnetDelete - Deletion of Subnet failed",
			Detail:   fmt.Sprintf("Deleting Subnet block by Id (%s) failed : %s", d.Id(), err),
		})
		return diags
	}
	// return empty string "" if deleted!
	d.SetId(refRes)
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
	return fmt.Sprintf("diamondip_subnet (ID = %s)", id)
}
