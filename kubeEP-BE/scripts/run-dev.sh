#!/usr/bin/env bash
set -e

go build -o app ./cmd/kubeEP && ./app