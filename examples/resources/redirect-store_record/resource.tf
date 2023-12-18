# Manage example record.
resource "redirect-store_record" "example" {
  name = "framework"
  to   = "https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework"
}
