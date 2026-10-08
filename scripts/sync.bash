#!/bin/bash
set -euxo pipefail

# 1. use the "oldstable" go toolchain.
releases=$(curl -fsS 'https://go.dev/dl/?mode=json')
oldstable=$(jq -r 'if length == 2 then .[1].version else "" end' <<<"$releases")
if [[ ! "$oldstable" =~ ^go(1\.[0-9]+)\.[0-9]+$ ]]; then
	echo "unexpected oldstable release: $oldstable"
	exit 1
fi
export GOTOOLCHAIN="$oldstable"

# 2. create temporary directory.
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

# 3. create temporary throw-away module.
(cd "$work" && go mod init example.com/bassosimone)

# 4. add each root dependency to the submodule trusting my local
# monorepo remote main branches as the source of truth.
for repo in "$@"; do
	git -C "$HOME/src/$repo" fetch --prune
	revision="$(git -C "$HOME/src/$repo" rev-parse remotes/origin/main)"
	(cd "$work" && go get -v "$repo@$revision")
done

# 5. start afresh with this module.
rm -rf pkg remotes
rm -f go.mod go.sum
go mod init github.com/bassosimone/nettree

# 6. install all non-bassosimone modules in a single command.
mods="$(cd "$work" && go list -m all | sed -n '2,$p' | sed -e 's/ /@/' | grep -v bassosimone)"
mapfile -t foreign_pins <<<"$mods"
go get -v "${foreign_pins[@]}"
git add go.mod go.sum

# 7. copy all bassosimone modules in tree.
mods="$(cd "$work" && go list -m all | sed -n '2,$p' | sed -e 's/ /@/' | grep bassosimone)"
mapfile -t owned_pins <<<"$mods"
for pin in "${owned_pins[@]}"; do
	orig_dir="$(go mod download -json "$pin" | jq -r .Dir)"
	path="$(go mod download -json "$pin" | jq -r .Path)"
	dirname="$(basename "$path")"
	rm -rf "pkg/$dirname"
	mkdir -p pkg
	cp -r "$orig_dir" "pkg/$dirname"
	chmod -R u+w "pkg/$dirname"
	mkdir -p remotes
	go mod download -json "$pin" |
		jq '{Path, Version, Sum, GoModSum, Origin}' \
			>"remotes/$dirname.json"
	./scripts/makereadme.bash "$dirname"
done

# 8. cleanup the directories content.
find pkg -type f -name go.mod -exec rm {} \;
find pkg -type f -name go.sum -exec rm {} \;
find pkg -depth -type d -name .github -exec rm -rf {} \;

# 9. add what survived.
git add ./pkg

# 10. rewrite the import paths.
find pkg -type f -name \*.go -exec \
	sed -i 's|"github.com/bassosimone/|"github.com/bassosimone/nettree/pkg/|' {} \;

# 11. update go.mod go.sum and import changes
go mod tidy
git add .
