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
	cp dist/linux/tdservice.service out/installer/tdservice.service
	cp dist/linux/install.sh out/installer/install.sh && chmod +x out/installer/install.sh
	cp out/tdservice out/installer/tdservice
	makeself --notemp out/installer out/tdservice-$(VERSION).bin "Threat Detection Service $(VERSION)" ./install.sh
	cp dist/linux/install_pgdb.sh out/install_pgdb.sh && chmod +x out/install_pgdb.sh

docker: installer
	cp dist/docker/entrypoint.sh out/entrypoint.sh && chmod +x out/entrypoint.sh
	docker build -t isecl/tdservice:latest -f ./dist/docker/Dockerfile ./out
	docker save isecl/tdservice:latest > ./out/docker-tdservice-$(VERSION)-$(GITCOMMIT).tar

docker-zip: installer
	mkdir -p out/docker-tdservice
	cp dist/docker/docker-compose.yml out/docker-tdservice/docker-compose
	cp dist/docker/entrypoint.sh out/docker-tdservice/entrypoint.sh && chmod +x out/docker-tdservice/entrypoint.sh
	cp dist/docker/README.md out/docker-tdservice/README.md
	cp out/tdservice-$(VERSION).bin out/docker-tdservice/tdservice-$(VERSION).bin
	cp dist/docker/Dockerfile out/docker-tdservice/Dockerfile
	zip -r out/docker-tdservice.zip out/docker-tdservice	

all: test docker

clean:
	rm -f cover.*
	rm -f tdservice
	rm -rf out/