---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "longship_organizational_units Data Source - terraform-provider-longship"
subcategory: ""
description: |-
  
---

# longship_organizational_units (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `organizational_units` (Attributes List) (see [below for nested schema](#nestedatt--organizational_units))

<a id="nestedatt--organizational_units"></a>
### Nested Schema for `organizational_units`

Read-Only:

- `address` (String)
- `city` (String)
- `code` (String)
- `company_email` (String)
- `country` (String)
- `customer_reference` (String)
- `direct_payment_profile_id` (String)
- `external_reference` (String)
- `financial_details` (Object) (see [below for nested schema](#nestedatt--organizational_units--financial_details))
- `grid_owner_reference` (String)
- `hotline_phone_number` (String)
- `house_number` (String)
- `id` (String)
- `msp_external_id` (String)
- `msp_ou_code` (String)
- `msp_ou_id` (String)
- `msp_ou_name` (String)
- `name` (String)
- `parent_id` (String)
- `postal_code` (String)
- `primary_contact_person` (String)
- `primary_contact_person_email` (String)
- `state` (String)
- `tenant_reference` (String)

<a id="nestedatt--organizational_units--financial_details"></a>
### Nested Schema for `organizational_units.financial_details`

Read-Only:

- `beneficiary_name` (String)
- `bic` (String)
- `iban` (String)
