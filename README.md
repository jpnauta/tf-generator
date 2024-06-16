## tf-generator

> Simple generator for Terraform/OpenTofu

## Getting Started

Simple generator program to combine various `.tf` and `.tfstate` files.

Unlike other DRY tools like [terragrunt](https://terragrunt.gruntwork.io/),
this tool does not add another layer on top of Terraform/OpenTofu. Instead, the
tool's sole purpose is to copy files into another generated file so it can be used
by multiple Terraform projects. This addresses Terraform's key limitation: it cannot
import files from other folders.

To see this tool in action, see the [examples/](./examples/) folder.
