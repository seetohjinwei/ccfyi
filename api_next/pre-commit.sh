#!/bin/sh

set -e

# change directory to this project
cd $(dirname "$0")

npm run test
