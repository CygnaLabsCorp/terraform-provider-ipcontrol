package main

import (
	"fmt"
	"terraform-provider-ipcontrol/ipcontrol"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = ipcontrol.Provider()
	testAccProviders = map[string]*schema.Provider{
		"ipcontrol": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := ipcontrol.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}

var server = fmt.Sprintf(
	`provider "ipcontrol" {
		server = "192.168.1.180"
		port = "1880"
		password = "incadmin"
		username = "incadmin"
	  }`)
