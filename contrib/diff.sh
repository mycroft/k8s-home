#!/bin/sh

git fetch origin
git checkout remotes/origin/generated -- generated/
diff -r generated/ dist/
git restore --staged generated
rm -fr generated/