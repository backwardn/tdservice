GITTAG := $(shell git describe --tags --abbrev=0 2> /dev/null)
GITCOMMIT := $(shell git describe --always)
GITCOMMITDATE := $(shell git log -1 --date=short --pretty=format:%cd)
VERSION := $(or ${GITTAG}, v0.0.0)

.PHONY: tdservice installer docker all test clean

tdservice:
	env GOOS=linux go build -ldflags "-X intel/isecl/tdservice/version.Version=$(version) -X intel/isecl/tdservice/version.GitHash=$(GITCOMMIT") -o out/tdservice

test:
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out
	go tool cover -html=cover.out -o cover.html

clean:
	rm -rf out/