terraform {
  required_providers {
    longship = {
      version = "0.1.2"
      source  = "cbcoutinho/longship"
    }
  }
}

provider "longship" {}
