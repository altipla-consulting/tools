#!/bin/bash

set -eu

. /opt/ci-toolset/functions.sh

GOOGLE_PROJECT=altipla-tools

configure-google-cloud

mkdir -p tmp/bin
run "GOBIN=tmp/bin go install ./cmd/..."
run "gsutil -h 'Cache-Control: no-cache' cp tmp/bin/* gs://tools.altipla.consulting/tools/"

run "gsutil -h 'Cache-Control: no-cache' cp install.sh gs://tools.altipla.consulting/install/tools"

git-tag
