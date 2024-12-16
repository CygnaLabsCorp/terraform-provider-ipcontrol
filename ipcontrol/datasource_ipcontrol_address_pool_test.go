package ipcontrol

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIPv4AddressPoolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my_addr_pool" {
							start_address       = "%s"
							end_address         = "%s"
							name                = "my-addrp"
							type                = "Dynamic DHCP"
							primary_net_service = "%s"
						}

						data "ipcontrol_address_pool" "my_pool" {
						  	start_address = "%s"
							depends_on = [ipcontrol_address_pool.my_addr_pool]
						}
						
					`, startIPv4AddrPoolTest, endIPv4AddrPoolTest, primaryNetServiceAddrPoolTest, startIPv4AddrPoolTest),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "container", "InControl/acctest"),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "start_address", startIPv4AddrPoolTest),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "end_address", endIPv4AddrPoolTest),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "name", "my-addrp"),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "type", "Dynamic DHCP"),
				),
			},
		},
	})
}

func TestAccIPv6AddressPoolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my_addr_pool_v6" {
							start_address       = "%s"
							prefix_length       = %v
							name                = "my-addrp-v6"
							type                = "Dynamic NA DHCPv6"
							primary_net_service = "%s"
						}

						data "ipcontrol_address_pool" "my_pool_v6" {
						  	start_address = "%s"
							depends_on = [ipcontrol_address_pool.my_addr_pool_v6]
						}
						
					`, startIPv6AddrPoolTest, prefixLength, primaryNetServiceIPv6AddrPoolTest, startIPv6AddrPoolTest),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "container", "InControl/acctest"),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "start_address", startIPv6AddrPoolTest),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "prefix_length", prefixLength),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "name", "my-addrp-v6"),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "type", "Dynamic NA DHCPv6"),
				),
			},
		},
	})
}
