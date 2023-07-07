# List all webhooks
data "longship_webhooks" "all" {}

output "webhooks" {
  value = data.longship_webhooks.all.webhooks
}
