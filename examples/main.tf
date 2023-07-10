terraform {
  required_providers {
    longship = {
      version = "0.1.12"
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

# --- Webhooks ---

data "longship_webhooks" "all" {
  depends_on = [
    longship_webhook.example
  ]
}

output "longship_webhooks" {
  value = data.longship_webhooks.all.webhooks
}

# --- Chargepoints ---

data "longship_chargepoints" "all" {}

output "longship_chargepoints" {
  value = data.longship_chargepoints.all.chargepoints
}

# --- Organizational Units ---

data "longship_organizational_units" "all" {}

output "longship_organizational_units" {
  value = data.longship_organizational_units.all.organizational_units
}
