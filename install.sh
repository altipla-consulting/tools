#!/bin/bash

set -eu

mkdir -p ~/bin

curl https://tools.altipla.consulting/tools/reloader > ~/bin/reloader
chmod +x ~/bin/reloader
