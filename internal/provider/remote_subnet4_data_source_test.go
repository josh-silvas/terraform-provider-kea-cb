package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccExampleDataSourceConfig = fmt.Sprintf(`
data "kea_remote_subnet4_data_source" "test" {
  hostname    = "%s"
  prefix      = "192.168.225.0/24"
}`, testAccHostname)

func TestAccRemoteSubnet4DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kea_remote_subnet4_data_source.test", "hostname", testAccHostname),
				),
			},
		},
	})
}
