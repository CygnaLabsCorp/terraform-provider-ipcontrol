package ipcontrol

import (
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccQipIPv4Subnet(t *testing.T) {
	resourceName := "cygnalabs_qip_ipv4_subnet.ipv4_subnet"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckQipIPv4SubnetDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{
				Config: testAccConfigWithProviderQIP(
					`
					resource "cygnalabs_qip_ipv4_subnet" "ipv4_subnet" {
					org_name= "Demo"
					subnet_address="10.0.0.0"
					subnet_mask = "255.255.255.0"
					network_address="10.0.0.0"
					warning_percent=90
					subnet_name="subnet_name"
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQipIPv4SubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_name", "Demo"),
					resource.TestCheckResourceAttr(resourceName, "subnet_address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "subnet_mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(resourceName, "network_address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "warning_percent", "90"),
					resource.TestCheckResourceAttr(resourceName, "subnet_name", "subnet_name"),
					resource.TestCheckResourceAttr(resourceName, "warning_type", "0"),
				),
			},
			// Step 2 update
			{
				Config: testAccConfigWithProviderQIP(
					`
					resource "cygnalabs_qip_ipv4_subnet" "ipv4_subnet" {
					org_name= "Demo"
					subnet_address="10.0.0.0"
					subnet_mask = "255.255.255.0"
					network_address="10.0.0.0"
					warning_percent=90
					subnet_name="new_name"
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQipIPv4SubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_name", "Demo"),
					resource.TestCheckResourceAttr(resourceName, "subnet_address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "subnet_mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(resourceName, "network_address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "warning_percent", "90"),
					resource.TestCheckResourceAttr(resourceName, "subnet_name", "new_name"),
					resource.TestCheckResourceAttr(resourceName, "warning_type", "0"),
				),
			},

			// step 3 Update subnet_mask

			{
				Config: testAccConfigWithProviderQIP(
					`
					resource "cygnalabs_qip_ipv4_subnet" "ipv4_subnet" {
					org_name= "Demo"
					subnet_address="10.0.0.0"
					subnet_mask = "255.255.255.128"
					network_address="10.0.0.0"
					warning_percent=90
					subnet_name="new_name"
					}`,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQipIPv4SubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_name", "Demo"),
					resource.TestCheckResourceAttr(resourceName, "subnet_address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "subnet_mask", "255.255.255.128"),
					resource.TestCheckResourceAttr(resourceName, "network_address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "warning_percent", "90"),
					resource.TestCheckResourceAttr(resourceName, "subnet_name", "new_name"),
					resource.TestCheckResourceAttr(resourceName, "warning_type", "0"),
				),
			},
		},
	})
}

// Helper function to check if subnet exists
func testAccCheckQipIPv4SubnetExists(n string) resource.TestCheckFunc {
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
			"addressVersion": "4",
		}

		_, err := objMgr.GetQipIPv4Subnet(query)
		if err != nil {
			return err
		}

		return nil
	}
}

// Helper function to check if subnet is destroyed
func testAccCheckQipIPv4SubnetDestroy(s *terraform.State) error {
	connector := testAccProvider.Meta().(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cygnalabs_qip_ipv4_subnet" {
			continue
		}

		// Construct query based on the resource's attributes
		query := map[string]string{
			"orgName":        rs.Primary.Attributes["org_name"],
			"subnetAddress":  rs.Primary.Attributes["subnet_address"],
			"addressVersion": "4",
		}

		_, err := objMgr.GetQipIPv4Subnet(query)
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
