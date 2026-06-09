# Terraform Provider Google Workspace
<a href="https://terraform.io">
    <img src="https://www.datocms-assets.com/2885/1620155116-brandhcterraformverticalcolor.svg" alt="Terraform logo" align="right" height="50" />
</a>

[![Releases](https://img.shields.io/github/release/SamuZad/terraform-provider-googleworkspace.svg)](https://github.com/SamuZad/terraform-provider-googleworkspace/releases)
[![LICENSE](https://img.shields.io/github/license/SamuZad/terraform-provider-googleworkspace.svg)](https://github.com/SamuZad/terraform-provider-googleworkspace/blob/main/LICENSE)
[![Unit tests](https://github.com/SamuZad/terraform-provider-googleworkspace/actions/workflows/test.yml/badge.svg)](https://github.com/SamuZad/terraform-provider-googleworkspace/actions/workflows/test.yml)

This Google Workspace provider for Terraform allows you to manage domains, users, and groups in your Google Workspace.

This is a community-maintained fork of the [`hashicorp/terraform-provider-googleworkspace`](https://github.com/hashicorp/terraform-provider-googleworkspace) provider, which was archived by HashiCorp on 2025-06-30. It is published to the Terraform Registry as [`SamuZad/googleworkspace`](https://registry.terraform.io/providers/SamuZad/googleworkspace/latest). Please [file issues](https://github.com/SamuZad/terraform-provider-googleworkspace/issues/new/choose) generously and detail your experience while using the provider. We welcome your feedback.

## Maintainers

This fork is maintained by [@SamuZad](https://github.com/SamuZad). It descends from the original `hashicorp/terraform-provider-googleworkspace`, which is no longer maintained upstream.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x (or [OpenTofu](https://opentofu.org))
-	[Go](https://golang.org/doc/install) >= 1.24 (to build the provider)

## Upgrading the provider

The Google Workspace provider doesn't upgrade automatically once you've started using it. After a new release you can run

```bash
terraform init -upgrade
```

to upgrade to the latest stable version of the Google Workspace provider. See the [Terraform website](https://www.terraform.io/docs/configuration/providers.html#provider-versions)
for more information on provider upgrades, and how to set version constraints on your provider.

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command or `make build`:
```sh
$ make build
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using The provider

See the [Google Workspace Provider documentation](https://registry.terraform.io/providers/SamuZad/googleworkspace/latest/docs) to get started using the
Google Workspace provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).
You can use [goenv](https://github.com/syndbg/goenv) to manage your Go version.
To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

For guidance on common development practices such as testing changes, see the [contribution guidelines](https://github.com/SamuZad/terraform-provider-googleworkspace/blob/main/.github/CONTRIBUTING.md).
If you have other development questions we don't cover, please file an issue!

## Special Recognition

* [Chase](https://github.com/DeviaVir) - for the excellent work creating the `DeviaVir/terraform-provider-gsuite` provider, the inspiration for this project.

## General Feedback
* How can we best support you ? - [feedback](https://forms.gle/XeqgPiFTtdevcRiu8)