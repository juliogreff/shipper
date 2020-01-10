#!/bin/bash -e

echo "::remove-matcher owner=go::"

E2E_FLAGS="--test.v" make setup e2e
