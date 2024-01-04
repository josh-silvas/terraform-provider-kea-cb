terraform {
  required_providers {
    kea = {
      source = "registry.terraform.io/josh-silvas/kea"
    }
  }
}

locals {
  hostname = "kea-primary.example.com"
  username = "admin"
  password = "password"
}

// Configure the Kea provider. This example relies on KEA_USERNAME and KEA_PASSWORD to exist
// in the environment. If not, you will need to specify them here.
provider "kea" {
#    username = local.username
#    password = local.password
}

data "kea_remote_subnet4_data_source" "example" {
  hostname = local.hostname
  prefix   = "192.168.168.0/24"
}


resource "kea_remote_subnet4_resource" "example_subnet4_resource" {
    hostname    = local.hostname
    subnet      = "192.168.168.0/24"
    pools       = [
      {pool = "192.168.168.50-192.168.168.150"}
    ]
    relay       = [
      {ip_address = "192.168.168.1"}
    ]
    option_data = [
      {code = 3,   name = "routers",             data = "192.168.168.1", always_send = true},
      {code = 15,  name = "domain-name",         data = "example.com", always_send = true},
      {code = 6,   name = "domain-name-servers", data = "4.2.2.2, 8.8.8.8", always_send = true},
    ]
}

resource "kea_reservation_resource" "example_reservation_resource" {
  hostname             = local.hostname
  reservation_hostname = "example-server"
  ip_address           = "192.168.168.160"
  subnet_id            = kea_remote_subnet4_resource.example_subnet4_resource.id
  hw_address           = "b8:27:eb:9c:ae:b7"
  option_data          = []
}
