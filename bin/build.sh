#!/bin/bash
set -ex

mkdir -p lambdapackage

cd $1

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o ../lambdapackage/collector

cp config.yml ../lambdapackage

cd ..

cp lambda/index.js lambdapackage

cd lambdapackage

zip $1.zip collector index.js config.yml

rm collector index.js config.yml
