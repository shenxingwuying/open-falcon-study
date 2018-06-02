#!/bin/bash

set -e

msg=$1

branch=$(git branch | grep "^*" | awk '{print $2}')

git commit -am "$msg"

git push origin $branch
