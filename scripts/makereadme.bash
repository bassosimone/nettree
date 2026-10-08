#!/bin/bash
set -euxo pipefail

dirname="$1"
commitHash="$(cat "./remotes/$dirname.json" | jq -r .Origin.Hash)"
repoURL="$(cat "./remotes/$dirname.json" | jq -r .Origin.URL)"
path="$(cat "./remotes/$dirname.json" | jq -r .Path)"
version="$(cat "./remotes/$dirname.json" | jq -r .Version)"

cat >"pkg/$dirname/README.md" <<EOF
# $dirname

Automatic import of $path into nettree.

- Upstream module: $path
- Version: $version
- GitHub tree: $repoURL/tree/$commitHash
- README.md: $repoURL/blob/$commitHash/README.md
- GitHub commit: $repoURL/commit/$commitHash

To verify:

\`\`\`bash
go mod download -json $path@$version
\`\`\`
EOF
