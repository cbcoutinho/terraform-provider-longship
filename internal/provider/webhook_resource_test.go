package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "longship_webhook" "test" {
  name = "test"
  ou_code = "0000"
  enabled = false
  event_types = ["SESSION_START"]
  url = "https://example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes of webhook
					resource.TestCheckResourceAttr("hashicups_order.test", "name", "test"),
					resource.TestCheckResourceAttr("hashicups_order.test", "ou_code", "0000"),
					resource.TestCheckResourceAttr("hashicups_order.test", "enabled", "false"),
					resource.TestCheckResourceAttr("hashicups_order.test", "event_types.0", "SESSION_START"),
					resource.TestCheckResourceAttr("hashicups_order.test", "url", "https://example.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("hashicups_order.test", "id"),
					resource.TestCheckResourceAttrSet("hashicups_order.test", "created"),
					resource.TestCheckResourceAttrSet("hashicups_order.test", "updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hashicups_order.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The created attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"created"},
			},
			{
				ResourceName:      "hashicups_order.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The created attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "hashicups_order" "test" {
  name = "test2"
  ou_code = "0000"
  enabled = false
  event_types = ["SESSION_START", "SESSION_STOP"]
  url = "https://example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("hashicups_order.test", "name", "test2"),
					resource.TestCheckResourceAttr("hashicups_order.test", "ou_code", "0000"),
					resource.TestCheckResourceAttr("hashicups_order.test", "enabled", "false"),
					resource.TestCheckResourceAttr("hashicups_order.test", "event_types.1", "SESSION_STOP"),
					resource.TestCheckResourceAttr("hashicups_order.test", "url", "https://example.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
