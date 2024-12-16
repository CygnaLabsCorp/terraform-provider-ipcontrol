package ipcontrol

import (
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	startIPv4AddrPoolTest         = "99.99.99.1"
	startIPv4AddrPoolTestUpdate   = "99.99.99.3"
	endIPv4AddrPoolTest           = "99.99.99.10"
	endIPv4AddrPoolTestUpdate     = "99.99.99.12"
	primaryNetServiceAddrPoolTest = "dhcp"

	startIPv6AddrPoolTest             = "2404:da1c:351:9000::"
	startIPv6AddrPoolTestUpdate       = "2404:da1c:351:9000:4d8b:cd2f:9128:6d00"
	prefixLength                      = "120"
	prefixLengthUpdate                = "121"
	primaryNetServiceIPv6AddrPoolTest = "dhcpv6"
)

func TestAccAddressPoolIPv4(t *testing.T) {
	resourceName := "ipcontrol_address_pool.my-addr-pool"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressPoolDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool" {
						start_address       = "%s"
						end_address         = "%s"
						name                = "my-addrp"
						type                = "Dynamic DHCP"
						primary_net_service = "%s"

						lifecycle {
							ignore_changes = [overlap_interface_ip, prefix_length]
						}
					}`, startIPv4AddrPoolTest, endIPv4AddrPoolTest, primaryNetServiceAddrPoolTest),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv4AddrPoolTest),
					resource.TestCheckResourceAttr(resourceName, "end_address", endIPv4AddrPoolTest),
					resource.TestCheckResourceAttr(resourceName, "type", "Dynamic DHCP"),
					resource.TestCheckResourceAttr(resourceName, "name", "my-addrp"),
				),
			},
			// Step 2 update start addr and end addr
			{
				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool" {
						start_address       = "%s"
						end_address         = "%s"
						name                = "my-addrp"
						type                = "Dynamic DHCP"
						primary_net_service = "%s"

						lifecycle {
							ignore_changes = [overlap_interface_ip, prefix_length]
						}
					}`, startIPv4AddrPoolTestUpdate, endIPv4AddrPoolTestUpdate, primaryNetServiceAddrPoolTest),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv4AddrPoolTestUpdate),
					resource.TestCheckResourceAttr(resourceName, "end_address", endIPv4AddrPoolTestUpdate),
					resource.TestCheckResourceAttr(resourceName, "type", "Dynamic DHCP"),
					resource.TestCheckResourceAttr(resourceName, "name", "my-addrp"),
				),
			},

			// step 3 Update name

			{
				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool" {
						start_address       = "%s"
						end_address         = "%s"
						name                = "my-addrp-update"
						type                = "Dynamic DHCP"
						primary_net_service = "%s"

						lifecycle {
							ignore_changes = [overlap_interface_ip, prefix_length]
						}
					}`, startIPv4AddrPoolTestUpdate, endIPv4AddrPoolTestUpdate, primaryNetServiceAddrPoolTest),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv4AddrPoolTestUpdate),
					resource.TestCheckResourceAttr(resourceName, "end_address", endIPv4AddrPoolTestUpdate),
					resource.TestCheckResourceAttr(resourceName, "type", "Dynamic DHCP"),
					resource.TestCheckResourceAttr(resourceName, "name", "my-addrp-update"),
				),
			},
		},
	})
}

func TestAccAddressPoolIPv6(t *testing.T) {
	resourceName := "ipcontrol_address_pool.my-addr-pool-v6"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressPoolDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool-v6" {
						start_address       = "%s"
						prefix_length       = %v
						name                = "my-addrp-v6"
						type                = "Dynamic NA DHCPv6"
						primary_net_service = "%s"

					}`, startIPv6AddrPoolTest, prefixLength, primaryNetServiceIPv6AddrPoolTest),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv6AddrPoolTest),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", prefixLength),
					resource.TestCheckResourceAttr(resourceName, "primary_net_service", primaryNetServiceIPv6AddrPoolTest),
					resource.TestCheckResourceAttr(resourceName, "type", "Dynamic NA DHCPv6"),
					resource.TestCheckResourceAttr(resourceName, "name", "my-addrp-v6"),
				),
			},
			// Step 2 update start addr and prefix length
			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool-v6" {
						start_address       = "%s"
						prefix_length       = %v
						name                = "my-addrp-v6"
						type                = "Dynamic NA DHCPv6"
						primary_net_service = "%s"

					}`, startIPv6AddrPoolTestUpdate, prefixLengthUpdate, primaryNetServiceIPv6AddrPoolTest),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv6AddrPoolTestUpdate),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", prefixLengthUpdate),
					resource.TestCheckResourceAttr(resourceName, "primary_net_service", primaryNetServiceIPv6AddrPoolTest),
					resource.TestCheckResourceAttr(resourceName, "type", "Dynamic NA DHCPv6"),
					resource.TestCheckResourceAttr(resourceName, "name", "my-addrp-v6"),
				),
			},

			// step 3 Update name , policy set , option set

			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool-v6" {
						start_address       = "%s"
						prefix_length       = %v
						name                = "my-addrp-v6-update"
						type                = "Dynamic NA DHCPv6"
						primary_net_service = "%s"
						dhcp_option_set      = "Cisco DHCPv6 Option Set"
						dhcp_policy_set      = "Cisco DHCP 8.0 Client Class Template Policy Set"

					}`, startIPv6AddrPoolTestUpdate, prefixLengthUpdate, primaryNetServiceIPv6AddrPoolTest),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv6AddrPoolTestUpdate),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", prefixLengthUpdate),
					resource.TestCheckResourceAttr(resourceName, "primary_net_service", primaryNetServiceIPv6AddrPoolTest),
					resource.TestCheckResourceAttr(resourceName, "type", "Dynamic NA DHCPv6"),
					resource.TestCheckResourceAttr(resourceName, "name", "my-addrp-v6-update"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_set", "Cisco DHCPv6 Option Set"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_policy_set", "Cisco DHCP 8.0 Client Class Template Policy Set"),
				),
			},
		},
	})
}

// Helper function to check if subnet exists
func testAccCheckAddressPoolExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		connector := testAccProvider.Meta().(*cc.Connector)
		objMgr := cc.NewObjectManager(connector)

		// Construct query based on the resource's attributes
		query := map[string]string{
			"startAddress": rs.Primary.Attributes["start_address"],
		}

		_, err := objMgr.GetAddressPool(query)
		if err != nil {
			return err
		}

		return nil
	}
}

// Helper function to check if subnet is destroyed
func testAccCheckAddressPoolDestroy(s *terraform.State) error {
	connector := testAccProvider.Meta().(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ipcontrol_address_pool" {
			continue
		}

		// Construct query based on the resource's attributes
		query := map[string]string{
			"address":      rs.Primary.Attributes["address"],
			"container":    rs.Primary.Attributes["container"],
			"size":         rs.Primary.Attributes["size"],
			"rawcontainer": rs.Primary.Attributes["rawcontainer"],
		}

		_, err := objMgr.GetSubnet(query)
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
