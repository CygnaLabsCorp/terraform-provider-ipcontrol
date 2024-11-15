package ipcontrol

import (
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccASubnet(t *testing.T) {
	resourceName := "cygnalabs_ipc_subnet.my-ipc-subnet-2"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{

				Config: testAccConfigWithProvider(
					`
					resource "cygnalabs_ipc_subnet" "my-ipc-subnet-2" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "13.0.0.0"
						size=24
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "address", "13.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "size", "24"),
				),
			},
			// Step 2 update
			{
				Config: testAccConfigWithProvider(
					`
					resource "cygnalabs_ipc_subnet" "my-ipc-subnet-2" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "13.0.0.0"
						size=24
						name = "update name test"
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "address", "13.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "size", "24"),
					resource.TestCheckResourceAttr(resourceName, "name", "update name test"),
				),
			},

			// step 3 Update size

			{
				Config: testAccConfigWithProvider(
					`
					resource "cygnalabs_ipc_subnet" "my-ipc-subnet-2" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "13.0.0.0"
						size = 25
						name = "update size test"
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "address", "13.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "size", "25"),
					resource.TestCheckResourceAttr(resourceName, "name", "update size test"),
				),
			},
		},
	})
}

func TestAccAAAASubnet(t *testing.T) {
	resourceName := "cygnalabs_ipc_subnet.my-ipc-subnet-v6"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{

				Config: testAccConfigWithProvider(
					`
					resource "cygnalabs_ipc_subnet" "my-ipc-subnet-v6" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "2a04:2880:10ff:8001::"
						address_version = 6
						size = 121
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "address", "2a04:2880:10ff:8001::"),
					resource.TestCheckResourceAttr(resourceName, "address_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "size", "121"),
				),
			},
			// Step 2 update name
			{
				Config: testAccConfigWithProvider(
					`
					resource "cygnalabs_ipc_subnet" "my-ipc-subnet-v6" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "2a04:2880:10ff:8001::"
						address_version = 6
						size = 121
						name = "update name for IPv6"
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "address", "2a04:2880:10ff:8001::"),
					resource.TestCheckResourceAttr(resourceName, "size", "121"),
					resource.TestCheckResourceAttr(resourceName, "name", "update name for IPv6"),
				),
			},

			// step 3 Update size

			{
				Config: testAccConfigWithProvider(
					`
					resource "cygnalabs_ipc_subnet" "my-ipc-subnet-v6" {
						rawcontainer = true
						container = "InControl/acctest"
						address = "2a04:2880:10ff:8001::"
						address_version = 6
						size = 122
						name = "update size for IPv6"
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container", "InControl/acctest"),
					resource.TestCheckResourceAttr(resourceName, "address", "2a04:2880:10ff:8001::"),
					resource.TestCheckResourceAttr(resourceName, "size", "122"),
					resource.TestCheckResourceAttr(resourceName, "address_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "name", "update size for IPv6"),
				),
			},
		},
	})
}

// Helper function to check if subnet exists
func testAccCheckSubnetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet ID is set")
		}

		connector := testAccProvider.Meta().(*cc.Connector)
		objMgr := cc.NewObjectManager(connector)

		// Construct query based on the resource's attributes
		query := map[string]string{
			"address":      rs.Primary.Attributes["address"],
			"container":    rs.Primary.Attributes["container"],
			"size":         rs.Primary.Attributes["size"],
			"rawcontainer": rs.Primary.Attributes["rawcontainer"],
		}

		_, err := objMgr.GetSubnet(query)
		if err != nil {
			return err
		}

		return nil
	}
}

// Helper function to check if subnet is destroyed
func testAccCheckSubnetDestroy(s *terraform.State) error {
	connector := testAccProvider.Meta().(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cygnalabs_ipc_subnet" {
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
