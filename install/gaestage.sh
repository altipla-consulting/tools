#!/bin/bash

set -eu

if ! command -v jq &> /dev/null
then
    echo "jq could not be found"
    echo "install it with `sudo apt install jq` before running this script"
    exit
fi

TOOL=gaestage
LATEST_VERSION=$(curl -sL https://api.github.com/repos/altipla-consulting/tools/releases/latest | jq -r '.tag_name')

echo "----------"
echo " [*] downloading $TOOL $LATEST_VERSION"
echo "----------"
curl -L https://github.com/altipla-consulting/tools/releases/download/${LATEST_VERSION}/${TOOL}_${LATEST_VERSION}_linux_amd64 > /usr/local/bin/$TOOL
chmod +x /usr/local/bin/$TOOL
