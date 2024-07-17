#!/usr/bin/env bash
set -e

nix-build -E 'let tf-generator-overlay = import ./.; pkgs = import <nixpkgs> { overlays = [ tf-generator-overlay ]; }; in pkgs.mkShellNoCC { buildInputs = with pkgs; [ tf-generator ]; }' --dry-run
