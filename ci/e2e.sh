#!/bin/bash -e

echo "::remove-matcher owner=go::"

E2E_FLAGS="--test.v --test.failfast" make setup e2e
