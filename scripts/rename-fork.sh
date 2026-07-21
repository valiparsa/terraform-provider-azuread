#!/usr/bin/env bash
#
# rename-fork.sh — rewrite upstream identifiers to this fork's.
#
# Re-run this after every upstream version bump (it is idempotent). It replaces
# the historical one-off "rename hashicorp to the forked repo" commit, so that
# files newly added by upstream are also converted.
#
# NOTE: this script excludes itself from the rewrite — it necessarily contains
# the literal upstream path in the sed patterns below, and must not clobber them.
#
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

self='scripts/rename-fork.sh'
old_mod='hashicorp/terraform-provider-azuread'
new_mod='valiparsa/terraform-provider-azuread'

# 1. Go module / import path (go.mod + all .go files), excluding this script.
# The grep exits non-zero when nothing matches (already renamed) — tolerate that
# so the script stays idempotent and reaches step 2.
git ls-files -z -- . ":(exclude)${self}" \
  | { xargs -0 grep -lZ "${old_mod}" 2>/dev/null || true; } \
  | xargs -0 --no-run-if-empty sed -i "s|${old_mod}|${new_mod}|g"

# 2. Provider registry address advertised by the provider binary (main.go).
sed -i 's|registry\.terraform\.io/hashicorp/azuread|registry.terraform.io/GPKbdZZb/forked-azuread|g' main.go

echo "rename-fork: done. Review 'git diff', then run 'go mod tidy' and 'go build ./...'."
