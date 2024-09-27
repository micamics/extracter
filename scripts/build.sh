#!/usr/bin/env bash

set -e
set -o errexit

cd cmd
go build -o extracter .
cd ..
mv cmd/extracter .