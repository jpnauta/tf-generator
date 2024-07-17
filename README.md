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

## Usage

To use this tool, first install it onto your machine:

```
curl https://raw.githubusercontent.com/jpnauta/tf-generator/development/install.sh | sudo bash
```

Next, create a sample project with the `init` command.

```
tf-generator init
```

Finally, follow the instructions to generate your first generated file.

```
tf-generator generate
```

### Nix Usage

To use this library with [nix](https://nixos.org/), import the repository as an overlay.

```
let
  tf-generator-overlay = import (builtins.fetchTarball https://github.com/jpnauta/tf-generator/archive/development.tar.gz);
  pkgs = import <nixpkgs> { overlays = [ tf-generator-overlay ]; };
in
pkgs.mkShellNoCC {
  buildInputs = with pkgs; [ tf-generator ];
}
```
