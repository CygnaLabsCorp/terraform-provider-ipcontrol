package ipcontrol

import (
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccQipIPv6Subnet(t *testing.T) {
	resourceName := "cygnalabs_qip_ipv6_subnet.ipv6_subnet"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckQipIPv6SubnetDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{
				Config: testAccConfigWithProviderQIP(
					`
					resource "cygnalabs_qip_ipv6_subnet" "ipv6_subnet" {
					org_name= "Demo"
					subnet_address="2000:0:0:10::"
					subnet_name="subnet_name"
					subnet_prefix_length = 60
					block_prefix_length = 48
					block_address="2000::"
					create_reverse_zone=true
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQipIPv6SubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_name", "Demo"),
					resource.TestCheckResourceAttr(resourceName, "subnet_address", "2000:0:0:10::"),
					resource.TestCheckResourceAttr(resourceName, "subnet_name", "subnet_name"),
					resource.TestCheckResourceAttr(resourceName, "subnet_prefix_length", "60"),
					resource.TestCheckResourceAttr(resourceName, "block_prefix_length", "48"),
					resource.TestCheckResourceAttr(resourceName, "block_address", "2000::"),
				),
			},
			// Step 2 update
			{
				Config: testAccConfigWithProviderQIP(
					`
					resource "cygnalabs_qip_ipv6_subnet" "ipv6_subnet" {
					org_name= "Demo"
					subnet_address="2000:0:0:10::"
					subnet_name="new_name"
					subnet_prefix_length = 60
					block_prefix_length = 48
					block_address="2000::"
					create_reverse_zone=true
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQipIPv6SubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_name", "Demo"),
					resource.TestCheckResourceAttr(resourceName, "subnet_address", "2000:0:0:10::"),
					resource.TestCheckResourceAttr(resourceName, "subnet_name", "new_name"),
					resource.TestCheckResourceAttr(resourceName, "subnet_prefix_length", "60"),
					resource.TestCheckResourceAttr(resourceName, "block_prefix_length", "48"),
					resource.TestCheckResourceAttr(resourceName, "block_address", "2000::"),
				),
			},

			// step 3 Update subnet_mask

			{
				Config: testAccConfigWithProviderQIP(
					`
					resource "cygnalabs_qip_ipv6_subnet" "ipv6_subnet" {
					org_name= "Demo"
					subnet_address="2000:0:0:10::"
					subnet_name="new_name"
					subnet_prefix_length = 62
					block_prefix_length = 48
					block_address="2000::"
					create_reverse_zone=true
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQipIPv6SubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_name", "Demo"),
					resource.TestCheckResourceAttr(resourceName, "subnet_address", "2000:0:0:10::"),
					resource.TestCheckResourceAttr(resourceName, "subnet_name", "new_name"),
					resource.TestCheckResourceAttr(resourceName, "subnet_prefix_length", "62"),
					resource.TestCheckResourceAttr(resourceName, "block_prefix_length", "48"),
					resource.TestCheckResourceAttr(resourceName, "block_address", "2000::"),
				),
			},
		},
	})
}

// Helper function to check if subnet exists
func testAccCheckQipIPv6SubnetExists(n string) resource.TestCheckFunc {
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
			"orgName":        rs.Primary.Attributes["org_name"],
			"subnetAddress":  rs.Primary.Attributes["subnet_address"],
			"prefixLength":   rs.Primary.Attributes["subnet_prefix_length"],
			"addressVersion": "6",
		}

		_, err := objMgr.GetQipIPv6Subnet(query)
		if err != nil {
			return err
		}

		return nil
	}
}

// Helper function to check if subnet is destroyed
func testAccCheckQipIPv6SubnetDestroy(s *terraform.State) error {
	connector := testAccProvider.Meta().(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cygnalabs_qip_ipv6_subnet" {
			continue
		}

		// Construct query based on the resource's attributes
		query := map[string]string{
			"orgName":        rs.Primary.Attributes["org_name"],
			"subnetAddress":  rs.Primary.Attributes["subnet_address"],
			"prefixLength":   rs.Primary.Attributes["subnet_prefix_length"],
			"addressVersion": "6",
		}

		_, err := objMgr.GetQipIPv6Subnet(query)
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
