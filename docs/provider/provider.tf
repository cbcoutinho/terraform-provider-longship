terraform {
  required_providers {
    longship = {
      version = "<no value>"
      source  = "cbcoutinho/longship"
    }
  }
}

provider "longship" {
  host            = var.longship_host
  tenant_key      = var.longship_tenant_key
  application_key = var.longship_application_key
}

resource "longship_webhook" "example" {
  name        = "test"
  ou_code     = "0000"
  enabled     = false
  event_types = ["SESSION_START"]
  url         = "https://example.com"
  headers = {
    hello = "world"
  }
}

data "longship_webhooks" "all" {
  depends_on = [
    longship_webhook.example
  ]
}

data "longship_chargepoints" "all" {}

output "longship_webhooks" {
  value = data.longship_webhooks.all.webhooks
}

output "longship_chargepoints" {
  value = data.longship_chargepoints.all.chargepoints
}
