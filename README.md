# experimental-terraform-redirect-store

https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `make build` command:

```shell
make build
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `tmp` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

```shell
make all
./tmp/api-server
make testacc
```

### Use the Provider locally

Write `~/.terraformrc`.

```
provider_installation {

  dev_overrides {
      "github.com/berquerant/redirect-store" = "/path/to/experimental-terraform-redirect-store/tmp"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Run API server.

``` shell
make all
./tmp/api-server
```

Apply example.

``` shell
cd examples/provider-install-verification
terraform apply
```

Accessing `http://localhost:8030/c/framework` will redirect you to `https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework`.
