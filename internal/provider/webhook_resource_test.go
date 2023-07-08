package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookResource(t *testing.T) {
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
					resource.TestCheckResourceAttr("longship_webhook.test", "name", "test"),
					resource.TestCheckResourceAttr("longship_webhook.test", "ou_code", "0000"),
					resource.TestCheckResourceAttr("longship_webhook.test", "enabled", "false"),
					resource.TestCheckResourceAttr("longship_webhook.test", "event_types.0", "SESSION_START"),
					resource.TestCheckResourceAttr("longship_webhook.test", "url", "https://example.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("longship_webhook.test", "id"),
					resource.TestCheckResourceAttrSet("longship_webhook.test", "created"),
					resource.TestCheckResourceAttrSet("longship_webhook.test", "updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "longship_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "longship_webhook" "test" {
  name = "test2"
  ou_code = "0000"
  enabled = false
  event_types = ["SESSION_START", "SESSION_STOP"]
  url = "https://example.com"
  headers = {
	hello = "world"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("longship_webhook.test", "name", "test2"),
					resource.TestCheckResourceAttr("longship_webhook.test", "ou_code", "0000"),
					resource.TestCheckResourceAttr("longship_webhook.test", "enabled", "false"),
					resource.TestCheckResourceAttr("longship_webhook.test", "event_types.1", "SESSION_STOP"),
					resource.TestCheckResourceAttr("longship_webhook.test", "url", "https://example.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("longship_webhook.test", "id"),
					resource.TestCheckResourceAttrSet("longship_webhook.test", "created"),
					resource.TestCheckResourceAttrSet("longship_webhook.test", "updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
