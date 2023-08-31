resource "kea_remote_subnet4_resource" "example" {
  hostname = "kea-primary.example.com"
  subnet   = "192.168.225.0/24"
  pools = [
    { pool = "192.168.225.50-192.168.225.150" }
  ]
  relay = [
    { ip_address = "192.168.225.1" }
  ]
  option_data = [
    { code = 3, name = "routers", data = "192.168.225.1" },
    { code = 15, name = "domain-name", data = "example.com" },
    { code = 6, name = "domain-name-servers", data = "4.2.2.2, 8.8.8.8", always_send = true },
  ]
  user_context = {
    "foo" = "bar"
  }
}
