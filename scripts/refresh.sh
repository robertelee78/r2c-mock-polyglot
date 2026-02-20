#!/usr/bin/env bash
# Restore the mock repo to its committed clean state.
#
# Pre-ADR-041 r2c mutates the scan target in place — it writes r2c_acquired_deps/,
# build outputs, .r2c/ workspaces, *.db, etc. into this tree. Run this between
# scans to get back to a pristine target.
#
#   git reset --hard  reverts tracked files r2c edited
#   git clean -fdx    deletes every untracked/ignored file r2c created
#                     (-x removes gitignored build artifacts too)
#
# NOTE: this discards ALL uncommitted changes, including source edits you have
# not committed. Commit anything you want to keep first.
set -euo pipefail
cd "$(dirname "$0")/.."
git reset -q --hard HEAD
git clean -qfdx -e '!.git'
echo "refreshed to $(git rev-parse --short HEAD) — tree clean:"
git status --short || true
