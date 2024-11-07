package ipcontrol

import (
	"context"
	"log"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CAA_SERVER", nil),
				Description: "Diamond IP CAA server IP address.",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CAA_USERNAME", nil),
				Description: "User to authenticate with Diamond IP CAA.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CAA_PASSWORD", nil),
				Description: "Password to authenticate with Diamond IP CAA.",
			},
			"port": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CAA_PORT", "1880"),
				Description: "Port number used for connection to Diamond IP CAA.",
			},

			"sslverify": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SSLVERIFY", "false"),
				Description: "If true, CAA client will verify SSL certificates.",
			},
			"connect_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONNECT_TIMEOUT", 60),
				Description: "Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"diamondip_subnet": resourceSubnet(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"diamondip_subnets": dataSourceSubnets(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	log.Println("Configure Diamond IP Provider ...")

	var seconds int64
	seconds = int64(d.Get("connect_timeout").(int))
	hostConfig := cc.HostConfig{
		Host:     d.Get("server").(string),
		Port:     d.Get("port").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	transportConfig := cc.TransportConfig{
		SslVerify:          d.Get("sslverify").(bool),
		HttpRequestTimeout: time.Duration(seconds),
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	requestBuilder := &cc.CaaRequestBuilder{}
	requestor := &cc.CaaHttpRequestor{}

	c, err := cc.NewConnector(hostConfig, transportConfig, requestBuilder, requestor)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
