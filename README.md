# Terraform Provider for Codefresh

The official Terraform and OpenTofu Provider for [Codefresh](https://codefresh.io/).

Terraform Registry: [registry.terraform.io/providers/codefresh-io/codefresh](https://registry.terraform.io/providers/codefresh-io/codefresh/latest)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) `1.x.x`  or [OpenTofu](https://github.com/opentofu/opentofu/releases/latest) `1.x.x`.

## Using the Provider

In `versions.tf`:

```terraform
terraform {
  required_providers {
    codefresh = {
      version = "x.y.z" # Optional but recommended; replace with latest semantic version
       source = "registry.terraform.io/codefresh-io/codefresh" # registry.terraform.io/ is optional for Terraform users, but required for OpenTofu users
    }
  }
}
```

You can also download and extract the provider binary (`terraform-provider-codefresh`) from [releases](https://github.com/codefresh-io/terraform-provider-codefresh/releases).

## Building the Provider Locally

```sh
make install
```

## [Provider Documentation](./docs)

The documentation is generated using [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs).

See: [CONTRIBUTING.md](./CONTRIBUTING.md#documentation)

## Provider Configuration:

See the [Provider Documentation](https://registry.terraform.io/providers/codefresh-io/codefresh/latest/docs#schema).

The provider requires a Codefresh API in order to authenticate to the Codefresh API. Generate the API key [here](https://g.codefresh.io/user/settings) and set the scopes [according to the resources you wish to create](https://codefresh.io/docs/docs/integrations/codefresh-api/#access-scopes). Note that some resource require platform admin permissions and hence can only be created for on-prem installations and not our SaaS offering.

The key can be set as an environment variable:

```bash
export CODEFRESH_API_KEY='xyz'
```

## Testing the Provider

**NOTE:** Acceptance tests create real resources, including admin resources (accounts, users) so make sure that `CODEFRESH_API_KEY` is set to a Codefresh installation and an account that you are ok with being modified.

```bash
make testacc
```

## OpenTofu Support

This provider supports [OpenTofu](https://opentofu.org/).

[Equivalence Testing](https://github.com/opentofu/equivalence-testing) is performed on the `examples/` directory in order to ensure that the provider behaves identically when used by either  `terraform` and `tofu` binaries.

The [OpenTofu Registry](https://registry.opentofu.org/) seems to be unpublished at time of writing. As of now, the provider is only published to the [Terraform Registry](https://registry.terraform.io/providers/codefresh-io/codefresh/latest).

## Contributors

<a href="https://github.com/codefresh-io/terraform-provider-codefresh/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=codefresh-io/terraform-provider-codefresh" />
</a>

## Acknowledgements

_This provider was initialized by [LightStep](https://lightstep.com/)_.

## License

Copyright 2023 Codefresh.

The Codefresh Provider is available under [MPL2.0 license](./LICENSE).
