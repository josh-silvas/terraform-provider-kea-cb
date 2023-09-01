resource "kea_reservation_resource" "example" {
  hostname             = "kea-primary.example.com"
  reservation_hostname = "test.example.com"
  ip_address           = "192.168.230.122"
  subnet_id            = 1921682300
  hw_address           = "94:8e:d3:db:d8:c5"
  option_data = [
    { code = 3, name = "routers", data = "192.168.230.1", always_send = false },
    { code = 15, name = "domain-name", data = "example.com", always_send = false },
    { code = 6, name = "domain-name-servers", data = "4.4.2.2, 8.8.8.8", always_send = true },
  ]
}
