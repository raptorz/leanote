#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
exec "$ROOT/bin/pearlnote" \
  -importPath github.com/pearlnote/pearlnote \
  -srcPath "$ROOT/runtime" \
  -runMode prod "$@"
