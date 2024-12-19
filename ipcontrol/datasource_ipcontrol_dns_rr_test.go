package ipcontrol

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTXTDnsRR(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
					resource "ipcontrol_dns_rr" "my_rr" {
					  owner                = "%s"
					  domain               = "%s"
					  resource_record_type = "%s"
					  data                 = "%s"
					  ttl    			   = "%s"
					  comment			   = "%s"
					}
					
					data "ipcontrol_dns_rr" "my_a_rr" {
					  owner                = "%s"
					  domain               = "%s"
					  resource_record_type = "%s"
					  data                 = "%s"
					  depends_on = [ipcontrol_dns_rr.my_rr]
					}
					`, owner, domain, AType, dataIPv4, ttl, comment, owner, domain, AType, dataIPv4),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "owner", owner),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "domain", domain),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "resource_record_type", AType),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "data", dataIPv4),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "ttl", ttl),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "comment", comment),
				),
			},
		},
	})
}

func TestAccDataSourceDnsTXTRR(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
					resource "ipcontrol_dns_rr" "my_rr" {
					  owner                = "%s"
					  domain               = "%s"
					  resource_record_type = "%s"
					  data                 = "%s"
					  ttl    			   = "%s"
					  comment			   = "%s"
					}
					
					data "ipcontrol_dns_rr" "my_a_rr" {
					  owner                = "%s"
					  domain               = "%s"
					  resource_record_type = "%s"
					  data                 = "%s"
					  depends_on = [ipcontrol_dns_rr.my_rr]
					}
					`, owner, domain, TXTType, dataText, ttl, comment, owner, domain, TXTType, dataText),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "owner", owner),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "domain", domain),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "resource_record_type", TXTType),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "data", dataText),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "ttl", ttl),
					resource.TestCheckResourceAttr("data.ipcontrol_dns_rr.my_a_rr", "comment", comment),
				),
			},
		},
	})
}
