package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccExampleResourceConfig = fmt.Sprintf(`
resource "kea_remote_subnet4_resource" "test" {
    hostname    = "%s"
    subnet      = "192.168.225.0/24"
    pools       = [
      {pool = "192.168.225.50-192.168.225.150"}
    ]
    relay       = [
      {ip_address = "192.168.225.1"}
    ]
    option_data = [
      {code = 3,   name = "routers",             data = "192.168.225.1"},
      {code = 15,  name = "domain-name",         data = "example.com"},
      {code = 6,   name = "domain-name-servers", data = "4.2.2.2, 8.8.8.8"},
    ]
}`, testAccHostname)

func TestAccRemoteSubnet4Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccExampleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kea_remote_subnet4_resource.test", "hostname", testAccHostname),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccExampleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kea_remote_subnet4_resource.test", "hostname", testAccHostname),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
