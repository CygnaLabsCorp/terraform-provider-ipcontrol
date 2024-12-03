package ipcontrol

import (
	"context"
	"fmt"
	"log"
	"strconv"
	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAddressRecordContext,
		ReadContext:   getAddressRecordContext,
		UpdateContext: updateAddressRecordContext,
		DeleteContext: deleteAddressRecordContext,

		Schema: map[string]*schema.Schema{
			"options": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"aliases": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"container": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the container that contains the device.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP Address of the Device.",
			},
			"address_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address type of this interface IP address",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Domain name of the Device.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The host name of the Device.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the Device.",
			},
			"device_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the device.",
			},
			"domain_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of the domain.",
			},
			"duid": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The DHCP Unique Identifier for IPv6.",
			},
			"interfaces": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_type": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							MaxItems:    1,
							Description: "The address type of this interface IP address.",
						},
						"container": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "The name of the container that contains the device.",
						},
						"ip_address": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							MaxItems:    1,
							Description: "The IP address of this interface.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of this interface.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The id of this interface.",
						},
						// "exclude_from_discovery": {
						// 	Type:     schema.TypeString,
						// 	Optional: true,
						// 	Computed: true,
						// },
						// "manufacturer": {
						// 	Type:     schema.TypeString,
						// 	Optional: true,
						// 	Computed: true,
						// },
						// "relay_agent_circuit_id": {
						// 	Type:     schema.TypeList,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Optional: true,
						// 	Computed: true,
						// },
						// "relay_agent_remote_id": {
						// 	Type:     schema.TypeList,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Optional: true,
						// },
						// "sequence": {
						// 	Type:     schema.TypeInt,
						// 	Optional: true,
						// 	Computed: true,
						// },
						// "virtual": {
						// 	Type:     schema.TypeList,
						// 	Elem:     &schema.Schema{Type: schema.TypeBool},
						// 	Optional: true,
						// 	Computed: true,
						// },
					},
				},
			},
			"resource_record_flag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func setIPCAddressResource(d *schema.ResourceData, address en.IPCAddress) {
	id := strconv.Itoa(address.ID)
	d.SetId(id)
	d.Set("aliases", address.Alias)
	d.Set("address_type", address.AddressType)
	// d.Set("id", id)
	d.Set("container", address.Container)
	d.Set("device_type", address.DeviceType)
	d.Set("domain_name", address.DomainName)
	d.Set("hostname", address.HostName)
	d.Set("ip_address", address.IpAddress)
	d.Set("resource_record_flag", address.ResourceRecordFlag)

	interfaces := []map[string]interface{}{}

	for _, intf := range address.Interfaces {
		// id := strconv.Itoa(intf.ID)
		interfaceItem := map[string]interface{}{
			"container":    intf.Container,
			"address_type": intf.AddressType,
			"ip_address":   intf.IpAddress,
			"id":           intf.ID,
			"name":         intf.Name,
		}
		interfaces = append(interfaces, interfaceItem)
	}
	d.Set("interfaces", interfaces)

}

func createAddressRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	// get params in schema
	deviceType := d.Get("device_type").(string)
	domainName := d.Get("domain_name").(string)
	hostname := d.Get("hostname").(string)
	interfaces := d.Get("interfaces").([]interface{})
	optionsList, optionsOk := d.Get("options").([]interface{})

	address := en.IPCAddressPost{
		DeviceType: deviceType,
		DomainName: domainName,
		HostName:   hostname,
		Interfaces: []en.AddressInterface{},
	}

	// handle options params
	if optionsOk {
		options, err := cc.ToStringSlice(optionsList)

		if err == nil {
			address.Options = options
		}
	}

	// hanle nested param [interfaces]
	for _, intf := range interfaces {
		var addressInterface en.AddressInterface
		iface := intf.(map[string]interface{})

		addressType, ok := iface["address_type"].([]interface{})

		// check if addressType not exist
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Can't read attribute address_type",
				Detail:   fmt.Sprintf("Can't read attribute address_type: (%s), intf (%s) ", addressType, iface),
			})
			return diags
		}
		addressTypeSlice, err := cc.ToStringSlice(addressType)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when format field address_type",
				Detail:   fmt.Sprintf("Error when format address_type: (%s)", err),
			})
			return diags
		}
		addressInterface.AddressType = addressTypeSlice

		ipAddress, ok := iface["ip_address"].([]interface{})
		// check if ipAddress not exist
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Can't read attribute ipAddress",
				Detail:   fmt.Sprintf("Can't read attribute ipAddress: (%s) ", ipAddress),
			})
			return diags
		}
		ipAddressSlice, err := cc.ToStringSlice(ipAddress)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when format attribute ip_address",
				Detail:   fmt.Sprintf("Error when format attribute ip_address: (%s)", err),
			})
			return diags
		}
		addressInterface.IpAddress = ipAddressSlice

		// check if name exist
		name, ok := iface["name"].(string)
		if ok {
			addressInterface.Name = name
		}

		address.Interfaces = append(address.Interfaces, addressInterface)

	}
	resp, err := objMgr.CreateAddress(&address)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when create device address",
			Detail:   fmt.Sprintf("Error when create device address: (%s) ", err),
		})
		return diags
	}

	setIPCAddressResource(d, *resp)
	return diags
}

func getAddressRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	ipAddress := d.Get("ip_address").(string)

	query := map[string]string{
		"ipAddress": ipAddress,
	}

	response, err := objMgr.GetAddress(query)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when getting device address",
			Detail:   fmt.Sprintf("Error when getting device address (%v) : %s", response, err.Error()),
		})
		return diags
	}

	setIPCAddressResource(d, *response)
	log.Printf("[DEBUG] %s: Completed reading device address", rsSubnetIdString(d))

	return diags
}

func updateAddressRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	deviceType := d.Get("device_type").(string)
	domainName := d.Get("domain_name").(string)
	hostname := d.Get("hostname").(string)
	interfaces := d.Get("interfaces").([]interface{})
	deviceIdString := d.Get("id").(string)
	deviceId, err := strconv.Atoi(deviceIdString)
	optionsList, optionsOk := d.Get("options").([]interface{})

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't read device id attribute",
			Detail:   fmt.Sprintf("Can't read device id attribute: (%s)", err),
		})
		return diags
	}

	address := en.NewAddressPost(en.IPCAddressPost{
		DeviceType: deviceType,
		DomainName: domainName,
		HostName:   hostname,
		Interfaces: []en.AddressInterface{},
		ID:         deviceId,
	})

	// handle options params
	if optionsOk {
		options, err := cc.ToStringSlice(optionsList)

		if err == nil {
			address.Options = options
		}
	}

	// hanle nested param [interfaces]
	for _, intf := range interfaces {
		var addressInterface en.AddressInterface
		iface := intf.(map[string]interface{})

		addressType, ok := iface["address_type"].([]interface{})

		// check if addressType not exist
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Can't read address_type attribute",
				Detail:   fmt.Sprintf("Getting address_type failed: (%s), intf (%s) ", addressType, iface),
			})
			return diags
		}
		addressTypeSlice, err := cc.ToStringSlice(addressType)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when format address_type",
				Detail:   fmt.Sprintf("Error when format address_type: (%s)", err),
			})
			return diags
		}
		addressInterface.AddressType = addressTypeSlice

		ipAddress, ok := iface["ip_address"].([]interface{})
		// check if ipAddress not exist
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Can't read ip_address attribute",
				Detail:   fmt.Sprintf("Can't read ip_address attribute: (%s) ", ipAddress),
			})
			return diags
		}
		ipAddressSlice, err := cc.ToStringSlice(ipAddress)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when format ip_address",
				Detail:   fmt.Sprintf("Error when format ip_address: (%s)", err),
			})
			return diags
		}
		addressInterface.IpAddress = ipAddressSlice

		// check if name exist
		name, ok := iface["name"].(string)
		if ok {
			addressInterface.Name = name
		}

		id, ok := iface["id"].(int)
		// check if ipAddress not exist
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Can't read device interface id attribute",
				Detail:   fmt.Sprintf("Can't read device interface id attribute: (%v) ", id),
			})
			return diags
		}
		addressInterface.ID = id
		address.Interfaces = append(address.Interfaces, addressInterface)
	}

	resp, err := objMgr.UpdateAddress(address)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when updating device address",
			Detail:   fmt.Sprintf("Error when updating device address: (%v) ", err),
		})
		return diags
	}
	setIPCAddressResource(d, *resp)

	return diags
}

func deleteAddressRecordContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)
	ipAddress := d.Get("ip_address").(string)

	_, err := objMgr.DeleteAddressRef(ipAddress)

	if err != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when deletting device address",
			Detail:   fmt.Sprintf("Error when deletting device address: %s", err.Error()),
		}
		diags = append(diags, diag)
		return diags
	}
	d.SetId("")

	return diags
}
