#!/usr/bin/env sh
set -e
set -x

echo "Installing tf-generator..."
RELEASES_URL="https://github.com/jpnauta/tf-generator/releases"
FILE_BASENAME="tf-generator"
INSTALL_DIR="/usr/bin"
LATEST="$(curl -s https://api.github.com/repos/jpnauta/tf-generator/releases/latest \
| grep "tag_name" \
| cut -d : -f 2 \
| tr -d \" \
| tr -d , \
| tr -d " ")"

test -z "$VERSION" && VERSION="$LATEST"

test -z "$VERSION" && {
	echo "Unable to get tf-generator version." >&2
	exit 1
}

TMP_DIR="$(mktemp -d)"
# shellcheck disable=SC2064
trap "rm -rf \"$TMP_DIR\"" EXIT INT TERM

OS="$(uname -s)"
ARCH="$(uname -m)"
test "$ARCH" = "aarch64" && ARCH="arm64"
TAR_FILE="${FILE_BASENAME}_${OS}_${ARCH}.tar.gz"
CHECKSUMS_FILE="${FILE_BASENAME}_$(echo "$VERSION" | tr -d v)_checksums.txt"

(
	cd "$TMP_DIR"
	echo "Downloading tf-generator $VERSION..."
	curl -sfLO "$RELEASES_URL/download/$VERSION/$TAR_FILE"
	curl -sfLO "$RELEASES_URL/download/$VERSION/$CHECKSUMS_FILE"
	echo "Verifying checksums..."
	sha256sum --ignore-missing --quiet --check "$CHECKSUMS_FILE"
)

tar -xf "$TMP_DIR/$TAR_FILE" -C "$TMP_DIR"
mv "$TMP_DIR/tf-generator" "$INSTALL_DIR"
chmod +x "$INSTALL_DIR/tf-generator"
echo tf-generator installed to "$INSTALL_DIR/tf-generator"