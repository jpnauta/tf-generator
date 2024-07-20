#!/usr/bin/env bash
set -e

mkdir -p dist
nix-build -E 'let tf-generator-overlay = import ./.; pkgs = import <nixpkgs> { overlays = [ tf-generator-overlay ]; }; in pkgs.mkShellNoCC { buildInputs = with pkgs; [ tf-generator ]; }' -o dist/release
rm -rf dist
