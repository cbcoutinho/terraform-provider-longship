terraform {
  required_providers {
    longship = {
      version = "0.1"
      source  = "milence.com/data-platform/longship"
    }
  }
}

provider "longship" {}
