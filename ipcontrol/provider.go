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
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CAA_SERVER", nil),
				Description: "CAA server IP address.",
			},
			"username_ipc": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USERNAME_IPC", nil),
				Description: "User to authenticate with IPC.",
			},
			"password_ipc": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PASSWORD_IPC", nil),
				Description: "Password to authenticate with IPC.",
			},
			"username_qip": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USERNAME_QIP", nil),
				Description: "User to authenticate with QIP.",
			},
			"password_qip": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PASSWORD_QIP", nil),
				Description: "Password to authenticate with QIP.",
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CAA_PORT", "1880"),
				Description: "Port number used for connection to CAA.",
			},
			"context": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONTEXT", "workflow"),
				Description: "Context of CAA.",
			},
			"sslverify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SSLVERIFY", "false"),
				Description: "If true, CAA client will verify SSL certificates.",
			},
			"connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONNECT_TIMEOUT", 60),
				Description: "Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cygnalabs_ipc_subnet":      resourceSubnet(),
			"cygnalabs_qip_ipv4_subnet": resourceQipIPv4Subnet(),
			"cygnalabs_qip_ipv6_subnet": resourceQipIPv6Subnet(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cygnalabs_ipc_subnet":      dataSourceSubnets(),
			"cygnalabs_qip_ipv4_subnet": dataSourceQipIPv4Subnet(),
			"cygnalabs_qip_ipv6_subnet": dataSourceQipIPv6Subnet(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	log.Println("Configure Diamond IP Provider ...")

	var seconds int64
	seconds = int64(d.Get("connect_timeout").(int))
	hostConfig := cc.HostConfig{
		Host:        d.Get("server").(string),
		Port:        d.Get("port").(string),
		Context:     d.Get("context").(string),
		UsernameIPC: d.Get("username_ipc").(string),
		PasswordIPC: d.Get("password_ipc").(string),
		UsernameQIP: d.Get("username_qip").(string),
		PasswordQIP: d.Get("password_qip").(string),
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
