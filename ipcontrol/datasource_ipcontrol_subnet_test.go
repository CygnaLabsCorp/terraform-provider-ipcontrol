package ipcontrol

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccADataSourceIPC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_subnet" "my-ipc-subnet" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "13.0.0.0"
						size=24
					}
						
					data "ipcontrol_subnet" "myds" {
						container = "InControl/acctest"
						address = "13.0.0.0"
						rawcontainer = true
						size=24
						depends_on = [ipcontrol_subnet.my-ipc-subnet]
					}
					
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_subnet.myds", "container", "InControl/acctest"),
					resource.TestCheckResourceAttr("data.ipcontrol_subnet.myds", "address", "13.0.0.0"),
					resource.TestCheckResourceAttr("data.ipcontrol_subnet.myds", "size", "24"),
				),
			},
		},
	})
}

func TestAccAAAADataSourceIPC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWithProviderIPC(
					`
					resource "ipcontrol_subnet" "my-ipc-subnet-v6" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "2a04:2880:10ff:8001::"
						address_version = 6
						size = 121
					}
						
					data "ipcontrol_subnet" "mydsv6" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "2a04:2880:10ff:8001::"
						address_version = 6
						size = 121
						depends_on = [ipcontrol_subnet.my-ipc-subnet-v6]
					}
					
					`,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ipcontrol_subnet.mydsv6", "container", "InControl/acctest"),
					resource.TestCheckResourceAttr("data.ipcontrol_subnet.mydsv6", "address", "2a04:2880:10ff:8001::"),
					resource.TestCheckResourceAttr("data.ipcontrol_subnet.mydsv6", "size", "121"),
					resource.TestCheckResourceAttr("data.ipcontrol_subnet.mydsv6", "address_version", "6"),
				),
			},
		},
	})
}
