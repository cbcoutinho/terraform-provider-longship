terraform {
  required_providers {
    longship = {
      version = "0.1"
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
  url         = "https://google.com"
  headers = [{
    name  = "hello"
    value = "world"
  }]
}
