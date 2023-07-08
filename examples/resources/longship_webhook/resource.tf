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