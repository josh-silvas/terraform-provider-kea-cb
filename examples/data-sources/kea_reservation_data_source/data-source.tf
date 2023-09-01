data "kea_reservation_data_source" "example" {
  hostname          = "kea-primary.example.com"
  ip_or_mac_address = "192.168.230.21"
  subnet_id         = 1921682300
}
