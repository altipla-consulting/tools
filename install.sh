#!/bin/bash

set -eu

TOOLS=(
  ci
  gendc
  impsort
  jnet
  linter
  previewer-netlify
  pub
  reloader
  wave
)
for app in "${TOOLS[@]}"
do
  echo
  echo "----------"
  echo " [*] downloading $app"
  echo "----------"
  curl https://tools.altipla.consulting/tools/$app > /usr/local/bin/$app
  chmod +x /usr/local/bin/$app
  echo "----------"
done
