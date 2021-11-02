#!/bin/bash

set -eu

if ! command -v jq &> /dev/null
then
    echo "jq could not be found"
    echo "install it with `sudo apt install jq` before running this script"
    exit
fi

LATEST_VERSION=$(curl -sL https://api.github.com/repos/altipla-consulting/tools/releases/latest | jq -r '.tag_name')

TOOLS=(
  ci
  gaestage
  gendc
  impsort
  jnet
  linter
  previewer-netlify
  pub
  releaser
  reloader
  wave
)
for app in "${TOOLS[@]}"
do
  echo
  echo "----------"
  echo " [*] downloading $app $LATEST_VERSION"
  echo "----------"
  curl -L https://github.com/altipla-consulting/tools/releases/download/${LATEST_VERSION}/${app}_${LATEST_VERSION}_linux_amd64 > /usr/local/bin/$app
  chmod +x /usr/local/bin/$app

  echo "----------"
done
