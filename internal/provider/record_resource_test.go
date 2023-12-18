package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create record
			{
				Config: providerConfig + `resource "redirect-store_record" "test0" {
  name = "test0-name"
  to = "test0-to"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("redirect-store_record.test0", "name", "test0-name"),
					resource.TestCheckResourceAttr("redirect-store_record.test0", "to", "test0-to"),
				),
			},
			// Import state
			{
				ResourceName:      "redirect-store_record.test0",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the RedirectStore
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update record
			{
				Config: providerConfig + `resource "redirect-store_record" "test0" {
  name = "test0-name"
  to = "test0-to-changed"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("redirect-store_record.test0", "name", "test0-name"),
					resource.TestCheckResourceAttr("redirect-store_record.test0", "to", "test0-to-changed"),
				),
			},
		},
	})
}
