package ipcontrol

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccADataSourceAddressIPC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_address" "my_ipc_device" {
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

					}
						
					data "ipcontrol_address" "my_device" {
						ip_address = "8.0.0.2"
						depends_on = [ipcontrol_address.my_ipc_device]
					}
					
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_address.my_device", "container", "InControl/acctest"),
					resource.TestCheckResourceAttr("data.ipcontrol_address.my_device", "interfaces.0.ip_address.0", "8.0.0.2"),
					resource.TestCheckResourceAttr("data.ipcontrol_address.my_device", "interfaces.0.address_type.0", "Static"),
					resource.TestCheckResourceAttr("data.ipcontrol_address.my_device", "domain_name", "com."),
				),
			},
		},
	})
}

func TestAccAAAADataSourceAddressIPC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_address" "my_ipc_device_v6" {
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

					}

					data "ipcontrol_address" "my_device" {
						ip_address = "2001:db8:1a2b:3c4d::1"
						depends_on = [ipcontrol_address.my_ipc_device_v6]
					}
					
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_address.my_device", "container", "InControl/acctest"),
					resource.TestCheckResourceAttr("data.ipcontrol_address.my_device", "interfaces.0.ip_address.0", "2001:db8:1a2b:3c4d::1"),
					resource.TestCheckResourceAttr("data.ipcontrol_address.my_device", "interfaces.0.address_type.0", "Static"),
					resource.TestCheckResourceAttr("data.ipcontrol_address.my_device", "domain_name", "com."),
				),
			},
		},
	})
}
