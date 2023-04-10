#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"
cat *.tf > complete/full.tf
terraform fmt complete/full.tf
