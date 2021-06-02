#!/bin/bash

set -eu

mkdir -p ~/bin

TOOLS=(
  ci
  impsort
  jnet
  linter
  previewer-netlify
  pub
  reloader
)
for app in "${TOOLS[@]}"
do
  echo
  echo "----------"
  echo " [*] downloading $app"
  echo "----------"
  curl https://tools.altipla.consulting/tools/$app > ~/bin/$app
  chmod +x ~/bin/$app
  echo "----------"
done
