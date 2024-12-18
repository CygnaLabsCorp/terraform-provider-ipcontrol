package ipcontrol

import (
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAAddress(t *testing.T) {
	resourceName := "ipcontrol_address.my-ipc-resource"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{

				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_address" "my-ipc-resource" {
						options = [
							"ignoreDupWarning",
							"resourceRecordFlag"
						]
						
						device_type = "PC"
						domain_name = "com."
						hostname    = "tfhost"

						interfaces {
							name         = "tfname"
							address_type = ["Static"]
							ip_address   = ["8.0.0.2"]
						}

					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "PC"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "com."),
					resource.TestCheckResourceAttr(resourceName, "hostname", "tfhost"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address.0", "8.0.0.2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "tfname"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.address_type.0", "Static"),
				),
			},
			// Step 2 update first inteface ip address and host
			{
				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_address" "my-ipc-resource" {
						options = [
							"ignoreDupWarning",
							"resourceRecordFlag"
						]

						device_type = "PC"
						domain_name = "com."
						hostname    = "tfhostupdate"

						interfaces {
							name         = "tfname"
							address_type = ["Static"]
							ip_address   = ["8.0.0.3"]
						}

						interfaces {
							name         = "tfname2"
							address_type = ["Static"]
							ip_address   = ["8.0.0.4"]
						}

					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "PC"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "com."),
					resource.TestCheckResourceAttr(resourceName, "hostname", "tfhostupdate"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address.0", "8.0.0.3"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address.0", "8.0.0.4"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "tfname"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.name", "tfname2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.address_type.0", "Static"),
				),
			},

			// step 3 update second interface ip address and name
			{
				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_address" "my-ipc-resource" {
						options = [
							"ignoreDupWarning",
							"resourceRecordFlag"
						]

						device_type = "PC"
						domain_name = "com."
						hostname    = "tfhostupdate"

						interfaces {
							name         = "tfname"
							address_type = ["Static"]
							ip_address   = ["8.0.0.3"]
						}
						interfaces {
							name         = "tfname2-update"
							address_type = ["Static"]
							ip_address   = ["8.0.0.5"]
						}

					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "PC"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "com."),
					resource.TestCheckResourceAttr(resourceName, "hostname", "tfhostupdate"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address.0", "8.0.0.3"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address.0", "8.0.0.5"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "tfname"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.name", "tfname2-update"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.address_type.0", "Static"),
				),
			},
		},
	})
}

func TestAccAAAAAddress(t *testing.T) {
	resourceName := "ipcontrol_address.my-address-v6"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{

				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_address" "my-address-v6" {
						options = [
							"ignoreDupWarning",
							"resourceRecordFlag"
						]
						
						device_type = "PC"
						domain_name = "com."
						hostname    = "tfhostv6"

						interfaces {
							name         = "tfname"
							address_type = ["Static"]
							ip_address   = ["2001:db8:1a2b:3c4d::1"]
						}

						interfaces {
							name         = "tfname2"
							address_type = ["Static"]
							ip_address   = ["2001:db8:1a2b:3c4d::3"]
						}

					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "PC"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "com."),
					resource.TestCheckResourceAttr(resourceName, "hostname", "tfhostv6"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address.0", "2001:db8:1a2b:3c4d::1"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address.0", "2001:db8:1a2b:3c4d::3"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "tfname"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.name", "tfname2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.address_type.0", "Static"),
				),
			},
			// Step 2 update ip and host name
			{
				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_address" "my-address-v6" {
						options = [
							"ignoreDupWarning",
							"resourceRecordFlag"
						]
						
						device_type = "PC"
						domain_name = "com."
						hostname    = "tfhostv6update"

						interfaces {
							name         = "tfname"
							address_type = ["Static"]
							ip_address   = ["2001:db8:1a2b:3c4d::2"]
						}

						interfaces {
							name         = "tfname2"
							address_type = ["Static"]
							ip_address   = ["2001:db8:1a2b:3c4d::3"]
						}

					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "PC"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "com."),
					resource.TestCheckResourceAttr(resourceName, "hostname", "tfhostv6update"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address.0", "2001:db8:1a2b:3c4d::2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address.0", "2001:db8:1a2b:3c4d::3"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "tfname"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.name", "tfname2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.address_type.0", "Static"),
				),
			},

			// Step 3 update second interface ip address and name
			{
				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_address" "my-address-v6" {
						options = [
							"ignoreDupWarning",
							"resourceRecordFlag"
						]
						
						device_type = "PC"
						domain_name = "com."
						hostname    = "tfhostv6update"

						interfaces {
							name         = "tfname"
							address_type = ["Static"]
							ip_address   = ["2001:db8:1a2b:3c4d::2"]
						}

						interfaces {
							name         = "tfname2update"
							address_type = ["Static"]
							ip_address   = ["2001:db8:1a2b:3c4d::4"]
						}

					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddressExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "PC"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "com."),
					resource.TestCheckResourceAttr(resourceName, "hostname", "tfhostv6update"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address.0", "2001:db8:1a2b:3c4d::2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address.0", "2001:db8:1a2b:3c4d::4"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.name", "tfname"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.name", "tfname2update"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.address_type.0", "Static"),
				),
			},
		},
	})
}

// Helper function to check if subnet exists
func testAccCheckAddressExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Address ID is set")
		}

		connector := testAccProvider.Meta().(*cc.Connector)
		objMgr := cc.NewObjectManager(connector)

		// Construct query based on the resource's attributes
		query := map[string]string{
			"ipAddress": rs.Primary.Attributes["ip_address"],
		}

		_, err := objMgr.GetAddress(query)
		if err != nil {
			return err
		}

		return nil
	}
}

// Helper function to check if subnet is destroyed
func testAccCheckAddressDestroy(s *terraform.State) error {
	connector := testAccProvider.Meta().(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ipcontrol_address" {
			continue
		}

		// Construct query based on the resource's attributes
		query := map[string]string{
			"address": rs.Primary.Attributes["ip_address"],
		}

		_, err := objMgr.GetSubnet(query)
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
