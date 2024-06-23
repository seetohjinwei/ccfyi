#!/bin/sh

set -e

# change directory to this project
cd $(dirname "$0")

# TODO: make script fail if make tidy cleaned stuff up (rationale: commit will still go through, because these changes aren't staged)
make tidy > /dev/null
make test > /dev/null
