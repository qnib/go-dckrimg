#!/bin/bash
set -ex 

go get -d
go build -o go-github_$(git describe --abbrev=0 --tags)_$(uname -s)
