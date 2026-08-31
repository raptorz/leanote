#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 4 ]]; then
  echo "usage: $0 <version> <goos> <goarch> <output-dir>" >&2
  exit 2
fi

version=${1#v}
target_os=$2
target_arch=$3
output_dir=$4
project_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
package_name="pearlnote-${target_os}-${target_arch}-v${version}"
stage_dir=$(mktemp -d)
revel_stage=$(mktemp -d)
cleanup() {
  chmod -R u+w "$stage_dir" "$revel_stage" 2>/dev/null || true
  rm -rf "$stage_dir" "$revel_stage"
}
trap cleanup EXIT

cd "$project_root"
go run ./app/cmd build . "$revel_stage" prod

binary_name=pearlnote
migration_name=pearlnote-migrate
if [[ $target_os == windows ]]; then
  binary_name=pearlnote.exe
  migration_name=pearlnote-migrate.exe
fi

mkdir -p "$stage_dir/pearlnote/bin" \
  "$stage_dir/pearlnote/runtime/github.com/revel/revel" \
  "$stage_dir/pearlnote/runtime/github.com/pearlnote"

CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" \
  go build -trimpath -ldflags="-s -w" \
  -o "$stage_dir/pearlnote/bin/$binary_name" ./app/tmp
CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" \
  go build -trimpath -ldflags="-s -w" \
  -o "$stage_dir/pearlnote/bin/$migration_name" ./tools/migration

cp -R app/views conf messages public mongodb_backup database docs \
  "$stage_dir/pearlnote/"
cp README.md "$stage_dir/pearlnote/"
rm -rf "$stage_dir/pearlnote/public/upload"
mkdir -p "$stage_dir/pearlnote/public/upload"

revel_dir=$(go list -m -f '{{.Dir}}' github.com/revel/revel)
cp -R "$revel_dir/conf" "$revel_dir/templates" \
  "$stage_dir/pearlnote/runtime/github.com/revel/revel/"

if [[ $target_os == windows ]]; then
  cp scripts/release/run.bat "$stage_dir/pearlnote/run.bat"
else
  ln -s ../../.. "$stage_dir/pearlnote/runtime/github.com/pearlnote/pearlnote"
  cp scripts/release/run.sh "$stage_dir/pearlnote/run.sh"
  chmod +x "$stage_dir/pearlnote/run.sh"
fi

mkdir -p "$output_dir"
if [[ $target_os == windows ]]; then
  (
    cd "$stage_dir"
    zip -qr "$output_dir/$package_name.zip" pearlnote
  )
else
  tar -C "$stage_dir" -czf "$output_dir/$package_name.tar.gz" pearlnote
fi
