#!/bin/bash

set -e

msg=$1

git commit -am "$msg"
git push origin master
