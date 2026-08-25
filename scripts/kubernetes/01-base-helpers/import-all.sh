#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Import all base helpers
#  Source this file:  source "$(dirname "$0")/../01-base-helpers/import-all.sh"
# --------------------------------------------------------------------------

_HELPERS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$_HELPERS_DIR/logger.sh"
source "$_HELPERS_DIR/install-apt.sh"
