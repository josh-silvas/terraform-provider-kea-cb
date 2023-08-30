data "kea_remote_subnet4_data_source" "example" {
  hostname = "kea-primary.example.com"
  prefix   = "192.168.230.0/24"
}
