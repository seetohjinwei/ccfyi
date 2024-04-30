#!/bin/sh

set -e

# change directory to this project
cd $(dirname "$0")

make test
