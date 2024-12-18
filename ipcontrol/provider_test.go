package ipcontrol

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cygnalabs": testAccProvider,
	}
}

func testProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}

func testAccPreCheck(t *testing.T) {
	fmt.Println("No precheck conditions are currently in place; all prechecks will pass.")
	return
}

var serverIPC = fmt.Sprintf(
	`provider "cygnalabs" {
		server = "127.0.0.1"
		port = "1880"
		password_ipc = "incadmin"
		username_ipc = "incadmin"
	  }`)

var serverQIP = fmt.Sprintf(
	`provider "cygnalabs" {
	server = "127.0.0.1"
	port = "1880"
	password_qip = "qipman"
	username_qip = "qipman"
	}`)
