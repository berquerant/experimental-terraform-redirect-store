terraform {
  required_providers {
    redirect-store = {
      source = "github.com/berquerant/redirect-store"
    }
  }
}

provider "redirect-store" {
  endpoint = "http://127.0.0.1:8030"
}

resource "redirect-store_record" "example" {
  name = "framework"
  to   = "https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework"
}

data "redirect-store_records" "example" {}

output "redirect_records" {
  value = data.redirect-store_records.example
}
