package ipcontrol

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccADataSourceQipIPv4Subnet(t *testing.T) {
	dataName := "data.cygnalabs_qip_ipv4_subnet.ipv4_subnet_data"
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderQIP(
					`
					resource "cygnalabs_qip_ipv4_subnet" "ipv4_subnet_resource" {
					org_name= "Demo"
					subnet_address="10.0.0.0"
					subnet_mask = "255.255.255.0"
					network_address="10.0.0.0"
					warning_percent=90
					subnet_name="subnet_name"
					}
					
					data "cygnalabs_qip_ipv4_subnet" "ipv4_subnet_data" {
					org_name= "Demo"
					subnet_address="10.0.0.0"
					depends_on = [cygnalabs_qip_ipv4_subnet.ipv4_subnet_resource]
					}
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "org_name", "Demo"),
					resource.TestCheckResourceAttr(dataName, "subnet_address", "10.0.0.0"),
					resource.TestCheckResourceAttr(dataName, "subnet_mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(dataName, "network_address", "10.0.0.0"),
					resource.TestCheckResourceAttr(dataName, "warning_percent", "90"),
					resource.TestCheckResourceAttr(dataName, "subnet_name", "subnet_name"),
					resource.TestCheckResourceAttr(dataName, "warning_type", "0"),
				),
			},
		},
	})
}
