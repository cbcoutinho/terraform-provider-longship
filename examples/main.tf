terraform {
  required_providers {
    longship = {
      version = "0.1.9"
      source  = "cbcoutinho/longship"
    }
  }
}

provider "longship" {}

resource "longship_webhook" "example" {
  name    = "test"
  ou_code = "0000"
  enabled = false
  event_types = [
    "SESSION_START",
    "SESSION_UPDATE",
    "SESSION_STOP",
    "OPERATIONAL_STATUS",
    "CONNECTIVITY_STATUS",
    "CHARGEPOINT_BOOTED",
    "CDR_CREATED",
  ]
  url = "https://example.com"
  headers = {
    hello = "world"
  }
}

data "longship_webhooks" "all" {}
data "longship_chargepoints" "all" {}

output "longship_webhooks" {
  value = data.longship_webhooks.all.webhooks
}

output "longship_chargepoints" {
  value = data.longship_chargepoints.all.chargepoints
}
