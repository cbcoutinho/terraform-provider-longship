---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "longship_webhooks Data Source - terraform-provider-longship"
subcategory: ""
description: |-
  Fetches the list of webhooks
---

# longship_webhooks (Data Source)

Fetches the list of webhooks

## Example Usage

```terraform
# List all webhooks
data "longship_webhooks" "all" {}

output "webhooks" {
  value = data.longship_webhooks.all.webhooks
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `webhooks` (Attributes List) (see [below for nested schema](#nestedatt--webhooks))

<a id="nestedatt--webhooks"></a>
### Nested Schema for `webhooks`

Read-Only:

- `created` (String) Timestamp of when webhook was created.
- `enabled` (Boolean) Webhook enabled or not.
- `event_types` (List of String) Notifications triggered with this webhook.
- `id` (String) Unique identifier of the webhook.
- `name` (String) Name of the webhook.
- `updated` (String) Timestamp of when webhook was last updated.