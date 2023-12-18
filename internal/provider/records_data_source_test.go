package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRecordsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "redirect-store_records" "test0" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.redirect-store_records.test0", "records.#", "0"),
				),
			},
			// Add record
			{
				Config: providerConfig + `resource "redirect-store_record" "test1" {
  name = "test1-name"
  to = "test1-to"
}`,
			},
			{
				Config: providerConfig + `data "redirect-store_records" "test2" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.redirect-store_records.test2", "records.#", "1"),
					resource.TestCheckResourceAttr("data.redirect-store_records.test2", "records.0.name", "test1-name"),
					resource.TestCheckResourceAttr("data.redirect-store_records.test2", "records.0.to", "test1-to"),
				),
			},
		},
	})
}
