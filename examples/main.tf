terraform {
  required_providers {
    longship = {
      version = "0.1.2"
      source  = "milence.com/data-platform/longship"
    }
  }
}

provider "longship" {}

resource "longship_webhook" "example" {
  name        = "test"
  ou_code     = "0000"
  enabled     = false
  event_types = ["SESSION_START"]
  url         = "https://example.com"
  headers = [{
    name  = "hello"
    value = "world"
  }]
}

data "longship_webhooks" "all" {}

output "longship_webhooks" {
  value = data.longship_webhooks.all.webhooks
}
