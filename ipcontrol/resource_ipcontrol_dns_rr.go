package ipcontrol

import (
	"context"
	"fmt"
	"strconv"

	en "terraform-provider-ipcontrol/ipcontrol/entities"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDnsRR() *schema.Resource {
	return &schema.Resource{
		CreateContext: createRRContext,
		ReadContext:   getRRContext,
		UpdateContext: updateRRContext,
		DeleteContext: deleteRRContext,

		Schema: map[string]*schema.Schema{
			"domain_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Domain type already defined to IPControl. If not specified, the “Default” domain type will be used.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name where resource records are to be added.",
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The owner field of the resource record.",
			},
			"resource_record_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Class of the resource record. Defaults to “IN”.",
			},
			"resource_record_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of resource record being imported. To import a type not predefined in IPControl the type must be prepended with “Other”.",
			},
			"ttl": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Time to live",
			},
			"data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The data portion of the resource record. The format is dependent on the type specified above.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment text associated with the resource record.",
			},
			"device_rec_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "When true, this indicates that the resource record is bound to a device. When false, this indicates that the resource record is associated with the domain only, and not a specific device.",
			},
		},
	}
}

func setIPCDnsRRResource(d *schema.ResourceData, rr *en.IPCDnsRR) {
	id := strconv.Itoa(rr.ID)
	d.SetId(id)
	d.Set("data", rr.Data)
	d.Set("domain_type", rr.DomainType)
	d.Set("owner", rr.Owner)
	d.Set("domain", rr.Domain)
	d.Set("ttl", rr.TTL)
	d.Set("comment", rr.Comment)
	d.Set("resource_record_type", rr.ResourceRecType)
	d.Set("resource_record_class", rr.ResourceRecClass)
	d.Set("device_rec_flag", rr.DeviceRecFlag)
}

func createRRContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	owner := d.Get("owner").(string)
	domainType := d.Get("domain_type").(string)
	resourceRecType := d.Get("resource_record_type").(string)
	resourceRecClass := d.Get("resource_record_class").(string)
	domain := d.Get("domain").(string)
	data := d.Get("data").(string)
	comment := d.Get("comment").(string)
	ttl := d.Get("ttl").(string)
	deviceRecFlag := d.Get("device_rec_flag").(bool)

	payload := en.NewDnsRRPost(en.IPCDnsRRPost{
		Domain:           domain,
		Owner:            owner,
		DomainType:       domainType,
		ResourceRecType:  resourceRecType,
		Data:             data,
		Comment:          comment,
		TTL:              ttl,
		ResourceRecClass: resourceRecClass,
		DeviceRecFlag:    deviceRecFlag,
	})

	err := objMgr.CreateDnsRR(payload)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when creating dns resource record",
			Detail:   fmt.Sprintf("Error when creating dns resource record: (%v)", err),
		})
		return diags
	}

	return getRRContext(ctx, d, m)
}

func getRRContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	owner := d.Get("owner").(string)
	resourceRecType := d.Get("resource_record_type").(string)
	domain := d.Get("domain").(string)
	data := d.Get("data").(string)

	query := map[string]string{
		"owner":      owner,
		"domainName": domain,
		"type":       resourceRecType,
		"rdata":      data,
	}

	rr, err := objMgr.GetDnsRR(query)

	if err != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when getting dns resource record",
			Detail:   fmt.Sprintf("Error when getting dns resource record: (%v)", err),
		}
		diags = append(diags, diag)
		return diags
	}

	setIPCDnsRRResource(d, rr)
	return diags
}

func updateRRContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	idStr := d.Id()
	id, err := strconv.Atoi(idStr)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't format ID",
			Detail:   fmt.Sprintf("Error when format id : (%v)", err),
		})
		return diags
	}

	owner := d.Get("owner").(string)
	domainType := d.Get("domain_type").(string)
	resourceRecType := d.Get("resource_record_type").(string)
	resourceRecClass := d.Get("resource_record_class").(string)
	domain := d.Get("domain").(string)
	data := d.Get("data").(string)
	comment := d.Get("comment").(string)
	ttl := d.Get("ttl").(string)
	deviceRecFlag := d.Get("device_rec_flag").(bool)

	payload := en.NewDnsRRPost(en.IPCDnsRRPost{
		ID:               id,
		Domain:           domain,
		Owner:            owner,
		DomainType:       domainType,
		ResourceRecType:  resourceRecType,
		Data:             data,
		Comment:          comment,
		TTL:              ttl,
		ResourceRecClass: resourceRecClass,
		DeviceRecFlag:    deviceRecFlag,
	})

	err = objMgr.UpdateDnsRR(payload)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when updating dns resource record",
			Detail:   fmt.Sprintf("Error when updating dns resource record: (%v)", err),
		})
		return diags
	}

	return getRRContext(ctx, d, m)
}

func deleteRRContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connector := m.(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	owner := d.Get("owner").(string)
	resourceRecType := d.Get("resource_record_type").(string)
	domain := d.Get("domain").(string)
	data := d.Get("data").(string)

	err := objMgr.DeleteDnsRR(owner, domain, resourceRecType, data)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error when deleting dns resource record",
			Detail:   fmt.Sprintf("Error when deleting dns resource record: (%v)", err),
		})
		return diags
	}

	return diags
}
