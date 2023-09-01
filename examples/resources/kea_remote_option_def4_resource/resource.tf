resource "kea_remote_option_def4_resource" "example" {
  hostname = "kea-primary.example.com"
  code     = 223
  space    = "dhcp4"
  type     = "string"
  name     = "custom-option"
}