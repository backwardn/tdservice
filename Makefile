GITTAG := $(shell git describe --tags --abbrev=0 2> /dev/null)
GITCOMMIT := $(shell git describe --always)
GITCOMMITDATE := $(shell git log -1 --date=short --pretty=format:%cd)
VERSION := $(or ${GITTAG}, v0.0.0)

.PHONY: tdservice installer docker all test clean

tdservice:
	env GOOS=linux go build -ldflags "-X intel/isecl/tdservice/version.Version=$(VERSION) -X intel/isecl/tdservice/version.GitHash=$(GITCOMMIT)" -o out/tdservice

test:
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out
	go tool cover -html=cover.out -o cover.html


installer: tdservice
	mkdir -p out/installer
	cp dist/linux/install.sh out/installer/install.sh && chmod +x out/installer/install.sh
	cp out/tdservice out/installer/tdservice
	makeself out/installer out/tdservice-$(VERSION)-$(GITCOMMIT).bin "Threat Detection Service $(VERSION)" ./install.sh

all: test installer


clean:
	rm cover.*
	rm tdservice
	rm -rf out/