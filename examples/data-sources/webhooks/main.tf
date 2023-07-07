terraform {
  required_providers {
    longship = {
      version = "0.1"
      source  = "milence.com/data-platform/longship"
    }
  }
}

provider "longship" {}

data "longship_webhooks" "all" {}

output "all_webhooks" {
  value = data.longship_webhooks.all.webhooks
}
