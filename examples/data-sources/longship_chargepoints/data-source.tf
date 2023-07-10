provider "longship" {}

data "longship_chargepoints" "all" {}

output "longship_chargepoints" {
  value = longship_webhooks.all.webhooks
}
