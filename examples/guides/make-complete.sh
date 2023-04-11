#!/usr/bin/env bash

# Pass guide sub-directory name as positional argument #1

cd "$(dirname "${BASH_SOURCE[0]}")"
cd $1
mkdir -p complete
cat *.tf > complete/full.tf
terraform fmt complete/full.tf
