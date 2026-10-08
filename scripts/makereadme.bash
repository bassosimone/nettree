#!/bin/bash
set -euxo pipefail
dirname="$1"
path="$(cat "./remotes/$dirname.json" | jq -r .Path)"
version="$(cat "./remotes/$dirname.json" | jq -r .Version)"
cat >"pkg/$dirname/README.md" <<EOF
# $dirname

Automatic import of $path at $version into github.com/bassosimone/nettree.
EOF
