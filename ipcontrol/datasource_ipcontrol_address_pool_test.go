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
							name                = "%s"
							type                = "%v"
							primary_net_service = "%s"
						}

						data "ipcontrol_address_pool" "my_pool" {
						  	start_address = "%s"
							depends_on = [ipcontrol_address_pool.my_addr_pool]
						}
						
					`, startIPv4AddrPool, endIPv4AddrPool, addrPoolName, typeIPv4, primaryNetServiceAddrPool, startIPv4AddrPool),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "container", "InControl/acctest"),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "start_address", startIPv4AddrPool),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "end_address", endIPv4AddrPool),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "name", addrPoolName),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool", "type", typeIPv4),
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
							name                = "%s"
							type                = "%s"
							primary_net_service = "%s"
						}

						data "ipcontrol_address_pool" "my_pool_v6" {
						  	start_address = "%s"
							depends_on = [ipcontrol_address_pool.my_addr_pool_v6]
						}
						
					`, startIPv6AddrPool, prefixLength, addrPoolName, typeIPv6, primaryNetServiceIPv6AddrPool, startIPv6AddrPool),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "container", "InControl/acctest"),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "start_address", startIPv6AddrPool),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "prefix_length", prefixLength),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "name", addrPoolName),
					resource.TestCheckResourceAttr("data.ipcontrol_address_pool.my_pool_v6", "type", typeIPv6),
				),
			},
		},
	})
}
