package ipcontrol

import (
	"fmt"
	cc "terraform-provider-ipcontrol/ipcontrol/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	domain         = "com."
	owner          = "caa"
	AType          = "A"
	TXTType        = "TXT"
	ttl            = "12345"
	comment        = "test acc tf"
	dataIPv4       = "10.0.0.1"
	dataIPv4Update = "10.0.0.2"
	dataText       = "txt text record"
	dataTextUpdate = "txt text record update"
)

func TestAccADnsRR(t *testing.T) {
	resourceName := "ipcontrol_dns_rr.my_rr"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
					resource "ipcontrol_dns_rr" "my_rr" {
					  owner                = "%s"
					  domain               = "%s"
					  resource_record_type = "%s"
					  data                 = "%s"
					}
					`, owner, domain, AType, dataIPv4),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRRExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", owner),
					resource.TestCheckResourceAttr(resourceName, "data", dataIPv4),
					resource.TestCheckResourceAttr(resourceName, "resource_record_type", AType),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
				),
			},

			// Step 2 update data

			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
					resource "ipcontrol_dns_rr" "my_rr" {
					  owner                = "%s"
					  domain               = "%s"
					  resource_record_type = "%s"
					  data                 = "%s"
					}
					`, owner, domain, AType, dataIPv4Update),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRRExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", owner),
					resource.TestCheckResourceAttr(resourceName, "data", dataIPv4Update),
					resource.TestCheckResourceAttr(resourceName, "resource_record_type", AType),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
				),
			},

			// Step 3 add ttl, comment

			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
					resource "ipcontrol_dns_rr" "my_rr" {
					  owner                = "%s"
					  domain               = "%s"
					  resource_record_type = "%s"
					  data                 = "%s"
					  ttl                  = "%s"
					  comment              = "%s"
					}
					`, owner, domain, AType, dataIPv4Update, ttl, comment),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRRExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", owner),
					resource.TestCheckResourceAttr(resourceName, "data", dataIPv4Update),
					resource.TestCheckResourceAttr(resourceName, "resource_record_type", AType),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
					resource.TestCheckResourceAttr(resourceName, "ttl", ttl),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccTXTDnsRR(t *testing.T) {
	resourceName := "ipcontrol_dns_rr.my_txt_rr"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressDestroy,
		Steps: []resource.TestStep{
			//  Step 1 create
			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
								resource "ipcontrol_dns_rr" "my_txt_rr" {
								  owner                = "%s"
								  domain               = "%s"
								  resource_record_type = "%s"
								  data                 = "%s"
								}
								`, owner, domain, TXTType, dataText),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRRExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", owner),
					resource.TestCheckResourceAttr(resourceName, "data", dataText),
					resource.TestCheckResourceAttr(resourceName, "resource_record_type", TXTType),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
				),
			},

			// Step 2 update data

			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
								resource "ipcontrol_dns_rr" "my_txt_rr" {
								  owner                = "%s"
								  domain               = "%s"
								  resource_record_type = "%s"
								  data                 = "%s"
								}
								`, owner, domain, TXTType, dataTextUpdate),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRRExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", owner),
					resource.TestCheckResourceAttr(resourceName, "data", dataTextUpdate),
					resource.TestCheckResourceAttr(resourceName, "resource_record_type", TXTType),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
				),
			},

			// Step 3 add ttl, comment

			{

				Config: testAccConfigWithProviderIPC(
					fmt.Sprintf(`
								resource "ipcontrol_dns_rr" "my_txt_rr" {
								  owner                = "%s"
								  domain               = "%s"
								  resource_record_type = "%s"
								  data                 = "%s"
								  ttl                  = "%s"
								  comment              = "%s"
								}
								`, owner, domain, TXTType, dataTextUpdate, ttl, comment),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRRExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "owner", owner),
					resource.TestCheckResourceAttr(resourceName, "data", dataTextUpdate),
					resource.TestCheckResourceAttr(resourceName, "resource_record_type", TXTType),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
					resource.TestCheckResourceAttr(resourceName, "ttl", ttl),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

// Helper function to check if subnet exists
func testAccCheckDnsRRExists(n string) resource.TestCheckFunc {
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
			"owner":      rs.Primary.Attributes["owner"],
			"domainName": rs.Primary.Attributes["domain"],
			"type":       rs.Primary.Attributes["resource_record_type"],
			"rdata":      rs.Primary.Attributes["data"],
		}

		_, err := objMgr.GetDnsRR(query)
		if err != nil {
			return err
		}

		return nil
	}
}

// Helper function to check if subnet is destroyed
func testAccCheckDnsRRDestroy(s *terraform.State) error {
	connector := testAccProvider.Meta().(*cc.Connector)
	objMgr := cc.NewObjectManager(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ipcontrol_dns_rr" {
			continue
		}

		// Construct query based on the resource's attributes
		query := map[string]string{
			"owner":      rs.Primary.Attributes["owner"],
			"domainName": rs.Primary.Attributes["domain"],
			"type":       rs.Primary.Attributes["resource_record_type"],
			"rdata":      rs.Primary.Attributes["data"],
		}

		_, err := objMgr.GetDnsRR(query)
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
