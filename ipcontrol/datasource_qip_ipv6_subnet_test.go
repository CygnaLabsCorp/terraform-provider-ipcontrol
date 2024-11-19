package ipcontrol

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccADataSourceQipIPv6Subnet(t *testing.T) {
	dataName := "data.cygnalabs_qip_ipv6_subnet.ipv6_subnet_data"
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderQIP(
					`
					resource "cygnalabs_qip_ipv6_subnet" "ipv6_subnet_resource" {
					org_name= "Demo"
					subnet_address="2000:0:0:10::"
					subnet_name="subnet_name"
					subnet_prefix_length = 60
					block_prefix_length = 48
					block_address="2000::"
					create_reverse_zone=true
					}
					
					data "cygnalabs_qip_ipv6_subnet" "ipv6_subnet_data" {
					org_name= "Demo"
					subnet_address="2000:0:0:10::"
					subnet_prefix_length = 60
					depends_on = [cygnalabs_qip_ipv6_subnet.ipv6_subnet_resource]
					}
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "org_name", "Demo"),
					resource.TestCheckResourceAttr(dataName, "subnet_address", "2000:0:0:10::"),
					resource.TestCheckResourceAttr(dataName, "subnet_name", "subnet_name"),
					resource.TestCheckResourceAttr(dataName, "subnet_prefix_length", "60"),
				),
			},
		},
	})
}
