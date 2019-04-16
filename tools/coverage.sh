#!/usr/bin/env bash

PKG_LIST=$(go list ./... | grep -v /vendor/)
for package in ${PKG_LIST}; do
  go test -covermode=count -coverprofile "cover/${package##*/}.cov" "$package" ;
done
tail -q -n +2 cover/*.cov >> cover/coverage.cov
go tool cover -func=cover/coverage.cov