#!/bin/bash

set -eu

mkdir -p ~/bin

TOOLS=(
  reloader
  jnet
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
