#!/bin/bash

T=$(mktemp)
find . -type f | xargs -n 1 file | egrep "ELF"| awk -F: '{print $1}'
(cat .gitignore; find . -type f | xargs -n 1 file | egrep "ELF" | awk -F: '{print $1}' | sed -e 's%^\./%%') | sort | uniq >$T
mv $T .gitignore
