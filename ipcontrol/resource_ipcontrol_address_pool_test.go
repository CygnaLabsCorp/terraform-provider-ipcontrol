package ipcontrol

import (
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	startIPv4AddrPool             = "99.99.99.1"
	startIPv4AddrPoolUpdate       = "99.99.99.3"
	endIPv4AddrPool               = "99.99.99.10"
	endIPv4AddrPoolUpdate         = "99.99.99.12"
	primaryNetServiceAddrPool     = "dhcp"
	typeIPv4                      = "Dynamic DHCP"
	typeIPv6                      = "Dynamic NA DHCPv6"
	addrPoolName                  = "my-addrp"
	addrPoolNameUpdate            = "my-addrp-update"
	addrPoolNameIPv6              = "my-addrp-v6"
	addrPoolNameIPv6Update        = "my-addrp-v6-update"
	startIPv6AddrPool             = "2404:da1c:351:9000::"
	startIPv6AddrPoolUpdate       = "2404:da1c:351:9000:4d8b:cd2f:9128:6d00"
	prefixLength                  = "120"
	prefixLengthUpdate            = "121"
	primaryNetServiceIPv6AddrPool = "dhcpv6"
	dhcpOptionSet                 = "Cisco DHCPv6 Option Set"
	dhcpPolicySet                 = "Cisco DHCP 8.0 Client Class Template Policy Set"
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
						name                = "%s"
						type                = "%s"
						primary_net_service = "%s"

						lifecycle {
							ignore_changes = [overlap_interface_ip, prefix_length]
						}
					}`, startIPv4AddrPool, endIPv4AddrPool, addrPoolName, typeIPv4, primaryNetServiceAddrPool),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv4AddrPool),
					resource.TestCheckResourceAttr(resourceName, "end_address", endIPv4AddrPool),
					resource.TestCheckResourceAttr(resourceName, "type", typeIPv4),
					resource.TestCheckResourceAttr(resourceName, "name", addrPoolName),
				),
			},
			// Step 2 update start addr and end addr
			{
				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool" {
						start_address       = "%s"
						end_address         = "%s"
						name                = "%s"
						type                = "%s"
						primary_net_service = "%s"

						lifecycle {
							ignore_changes = [overlap_interface_ip, prefix_length]
						}
					}`, startIPv4AddrPoolUpdate, endIPv4AddrPoolUpdate, addrPoolName, typeIPv4, primaryNetServiceAddrPool),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv4AddrPoolUpdate),
					resource.TestCheckResourceAttr(resourceName, "end_address", endIPv4AddrPoolUpdate),
					resource.TestCheckResourceAttr(resourceName, "type", typeIPv4),
					resource.TestCheckResourceAttr(resourceName, "name", addrPoolName),
				),
			},

			// step 3 Update name

			{
				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool" {
						start_address       = "%s"
						end_address         = "%s"
						name                = "%s"
						type                = "%s"
						primary_net_service = "%s"

						lifecycle {
							ignore_changes = [overlap_interface_ip, prefix_length]
						}
					}`, startIPv4AddrPoolUpdate, endIPv4AddrPoolUpdate, addrPoolNameUpdate, typeIPv4, primaryNetServiceAddrPool),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv4AddrPoolUpdate),
					resource.TestCheckResourceAttr(resourceName, "end_address", endIPv4AddrPoolUpdate),
					resource.TestCheckResourceAttr(resourceName, "type", typeIPv4),
					resource.TestCheckResourceAttr(resourceName, "name", addrPoolNameUpdate),
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
						name                = "%s"
						type                = "%s"
						primary_net_service = "%s"

					}`, startIPv6AddrPool, prefixLength, addrPoolNameIPv6, typeIPv6, primaryNetServiceIPv6AddrPool),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv6AddrPool),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", prefixLength),
					resource.TestCheckResourceAttr(resourceName, "primary_net_service", primaryNetServiceIPv6AddrPool),
					resource.TestCheckResourceAttr(resourceName, "type", typeIPv6),
					resource.TestCheckResourceAttr(resourceName, "name", addrPoolNameIPv6),
				),
			},
			// Step 2 update start addr and prefix length
			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool-v6" {
						start_address       = "%s"
						prefix_length       = %v
						name                = "%s"
						type                = "%v"
						primary_net_service = "%s"

					}`, startIPv6AddrPoolUpdate, prefixLengthUpdate, addrPoolNameIPv6, typeIPv6, primaryNetServiceIPv6AddrPool),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv6AddrPoolUpdate),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", prefixLengthUpdate),
					resource.TestCheckResourceAttr(resourceName, "primary_net_service", primaryNetServiceIPv6AddrPool),
					resource.TestCheckResourceAttr(resourceName, "type", typeIPv6),
					resource.TestCheckResourceAttr(resourceName, "name", addrPoolNameIPv6),
				),
			},

			// step 3 Update name , policy set , option set

			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
						resource "ipcontrol_address_pool" "my-addr-pool-v6" {
						start_address       = "%s"
						prefix_length       = %v
						name                = "%s"
						type                = "%s"
						primary_net_service = "%s"
						dhcp_option_set      = "%s"
						dhcp_policy_set      = "%s"

					}`, startIPv6AddrPoolUpdate, prefixLengthUpdate, addrPoolNameIPv6Update, typeIPv6, primaryNetServiceIPv6AddrPool, dhcpOptionSet, dhcpPolicySet),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "start_address", startIPv6AddrPoolUpdate),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", prefixLengthUpdate),
					resource.TestCheckResourceAttr(resourceName, "primary_net_service", primaryNetServiceIPv6AddrPool),
					resource.TestCheckResourceAttr(resourceName, "type", typeIPv6),
					resource.TestCheckResourceAttr(resourceName, "name", addrPoolNameIPv6Update),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_set", dhcpOptionSet),
					resource.TestCheckResourceAttr(resourceName, "dhcp_policy_set", dhcpPolicySet),
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
