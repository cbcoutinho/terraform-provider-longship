terraform {
  required_providers {
    longship = {
      version = "0.1.9"
      source  = "cbcoutinho/longship"
    }
  }
}

provider "longship" {}
