#!/bin/bash -e

TEST_FLAGS="-v" make test lint verify-codegen
