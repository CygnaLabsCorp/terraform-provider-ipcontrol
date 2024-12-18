package ipcontrol

import (
	"context"
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDnsRRs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnsRRRead,
		Schema: map[string]*schema.Schema{
			"domain_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the domain type to which the domain belongs.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name of the resource records.",
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The owner field of the resource record.",
			},
			"resource_record_class": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Class of the resource record.",
			},
			"resource_record_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of resource record.",
			},
			"ttl": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time to live",
			},
			"data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The data portion of the resource record. The format is dependent on the type specified above.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment text associated with the resource record.",
			},
			"device_rec_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When true, this indicates that the resource record is bound to a device. When false, this indicates that the resource record is associated with the domain only, and not a specific device.",
			},
		},
	}
}

func dataSourceDnsRRRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
			Summary:  "Error when reading dns resource record",
			Detail:   fmt.Sprintf("Error when reading dns resource record: (%v)", err),
		}
		diags = append(diags, diag)
		return diags
	}

	setIPCDnsRRResource(d, rr)
	return diags
}
