# Terraform Provider Kea Configuration Backend

_This Terraform provider is designed to interact with [ISC Kea's Configuration Backend](https://kea.readthedocs.io/en/latest/arm/config.html#kea-configuration-backend)
to manage Kea configuration stored in a remote backend such as MySQL or PostgreSQL._

> **Note:** This provider requires the [libdhcp_db_cmds.so](https://kea.readthedocs.io/en/latest/arm/hooks.html#libdhcp-cb-cmds-so-configuration-backend-commands).
> hook library to be installed and configured on the Kea server. See the [Kea documentation](https://kea.readthedocs.io/en/latest/arm/) for more information.

This repository is built on Terraform scaffolding for providers and contains the following:

- A resource and a data source (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

# Table of Contents
1. [Requirements](#requirements)
2. [Building The Provider](#building-the-provider)
3. [Adding Dependencies](#adding-dependencies)
4. [Using The Provider](#using-the-provider)
5. [Developing The Provider](#developing-the-provider)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command or the Makefile `build` target.

```shell
make build
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get -u github.com/author/dependency
go mod tidy && go mod vendor
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

### Provider Configuration

```hcl
provider "kea" {
  username = "some-kea-ctrl-user"
  password = "some-kea-ctrl-password"
}
```

### Data Source Configuration
#### Remote Subet4 Commands
kea_remote_subnet4_data_source
```hcl
data "kea_remote_subnet4_data_source" "example" {
  hostname = "kea-primary.example.com"
  prefix   = "192.168.230.0/24"
}

```

### Resource Configuration
#### Remote Subet4 Commands
kea_remote_subnet4_resource
```hcl
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
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org)
installed on your machine (see [Requirements](#requirements) above).

See the [DEVELOPMENT](develop/README.md) documentation for more information.


## Resources

* [Josh Silvas](mailto:josh@jsilvas.com)

