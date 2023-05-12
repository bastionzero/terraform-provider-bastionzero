#!/usr/bin/env bash

# Find all links to BastionZero docs website and paste in .docs-links.txt file
echo "$(grep -h -r --only-matching 'docs.bastionzero.com[^) ]*' templates | sort --unique)" >| ./.docs-links.txt