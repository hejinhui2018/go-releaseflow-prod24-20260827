#!/bin/sh
set -eu
image="${1:-releaseflow:local}"
docker build -f benzhi.Dockerfile -t "$image" .

