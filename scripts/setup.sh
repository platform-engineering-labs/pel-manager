#!/usr/bin/env bash

set -e
set -o pipefail

export OS=$(uname | tr '[:upper:]' '[:lower:]')
export ARCH=$(uname -m |  tr -d '_' | sed s/aarch64/arm64/)

if ! which curl > /dev/null; then
  echo "curl not found in PATH, please install to continue"
  exit 1
fi

echo "Downloading: pelmgr"
if ! curl "https://hub.platform.engineering/get/binaries/${OS}-${ARCH}/pelmgr" 2>/dev/null > ./pelmgr; then
  echo "Failed to download: pelmgr"
  exit 1
fi

chmod +x ./pelmgr
./pelmgr "${@}"
rm -rf ./pelmgr
