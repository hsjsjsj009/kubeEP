#!/usr/bin/env bash
set -e

go build -o cron ./cmd/kubeEP-cron && ./cron