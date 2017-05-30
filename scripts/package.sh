#!/bin/sh

set -e
set -x

VERSION=$(grep "version = " cmd/root.go | sed -E 's/.*"(.+)"$/\1/')
REPO="pet"

rm -rf ./out/
gox --osarch "windows/386 windows/amd64 darwin/386 darwin/amd64 linux/386 linux/amd64" -output="./out/${REPO}_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}"

rm -rf ./pkg/
mkdir ./pkg

for PLATFORM in $(find ./out -mindepth 1 -maxdepth 1 -type d); do
        PLATFORM_NAME=$(basename ${PLATFORM})

        pushd ${PLATFORM}
        cp -r ../../misc ./
        zip -r ../../pkg/${PLATFORM_NAME}.zip ./*
        popd
done
