# Longship.io Terraform Provider

This provider enables automating EV Charging resources for the CPO [Longship](https://longship.io)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Here's a small example of how to use the provider:

```terraform
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
  event_types = ["SESSION_START"]
  url = "https://example.com"
  headers = {
    hello = "world"
  }
}

data "longship_webhooks" "all" {
  depends_on = [
    longship_webhook.example
  ]
}

data "longship_chargepoints" "all" {}

output "longship_webhooks" {
  value = data.longship_webhooks.all.webhooks
}

output "longship_chargepoints" {
  value = data.longship_chargepoints.all.chargepoints
}
```


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
