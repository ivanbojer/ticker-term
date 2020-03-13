SHELL=bash
.ONESHELL:

MAKEFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
SCRIPTPATH=$(shell dirname "$(MAKEFILE_PATH)")

PRIMARY_REPO=$(shell git remote -v | head -n 1)
REPO_NAME=$(shell echo "$(PRIMARY_REPO)" | sed -r 's_.*git@github.com:.*/([a-z0-9_-]*).*_\1_i')
REPO_PATH=$(shell echo "$(SCRIPTPATH)" | sed -r 's;(.*$(REPO_NAME)).*;\1;g')
SSH_URL=$(shell echo "$(PRIMARY_REPO)" | sed -r 's_.*(git@.*git).*_\1_')
GO_REPO=$(shell echo "$(SSH_URL)" | cut -d @ -f 2 | tr : / | cut -d . -f 1-2)
GO_PKG_PATH=$(shell echo "$(SCRIPTPATH)" | sed -r 's;.*$(REPO_NAME)(.*);$(GO_REPO)\1;g')
GO_PKG_NAME=$(shell echo "$(GO_PKG_PATH)" | rev | cut -d / -f 1 | rev)
BRANCH_NAME=$(shell git branch | grep '*' | tr -d ' *')
COMMIT_SHA=$(shell git show HEAD --format=%H | head -n 1)

GIT_CURRENT_TAG=$(shell git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///')
GIT_LAST_TAG=$(shell git tag | tail -n 1)
GIT_TAG=$(shell echo $(GIT_LAST_TAG) | tr -d v )

COMMIT_SHA_SHORT=$(shell git rev-parse --short HEAD)
TIMESTAMP_UNIX=$(shell date --utc +%s)
TIMESTAMP_GOMOD=$(shell date --utc +%Y%m%d%H%M%S)
DATE_WITH_TIMESTAMP=$(shell date --utc +%d%m%y-%s)

ARGS=$(filter-out $@,$(MAKECMDGOALS))
1=$(shell echo $(ARGS) | cut -d " " -f 1)
2=$(shell echo $(ARGS) | cut -d " " -f 2)
3=$(shell echo $(ARGS) | cut -d " " -f 3)

.PHONY: $(ARGS)
.SILENT:


MAKEFLAGS += --no-print-directory
ARGS=$(filter-out $@,$(MAKECMDGOALS))
1=$(shell echo $(ARGS) | cut -d " " -f 1)
2=$(shell echo $(ARGS) | cut -d " " -f 2)
3=$(shell echo $(ARGS) | cut -d " " -f 3)

.PHONY: $(ARGS)
.SILENT:


MAKEFLAGS += --no-print-directory
##############################################
all: build.linux build.darwin

build.linux: 
	mkdir -p $(SCRIPTPATH)/build/linux
	GOSUMDB=off GOPROXY=direct CGO_ENABLED=0 GOOS=linux \
		go build -a \
		--ldflags "-s -w \
		-extldflags \"-static\"" \
		-installsuffix cgo \
	    -o $(SCRIPTPATH)/build/linux/bin/ticker-term

	tar -cf \
		$(SCRIPTPATH)/build/linux/ticker-term-$(GIT_TAG)-linux.tar.gz \
		-C $(SCRIPTPATH)/build/linux/bin .
	
build.darwin: 
	mkdir -p $(SCRIPTPATH)/build/darwin
	GOSUMDB=off GOPROXY=direct CGO_ENABLED=0 GOOS=darwin \
		go build -a \
		--ldflags "-s -w \
		-extldflags \"-static\"" \
		-installsuffix cgo \
	    -o $(SCRIPTPATH)/build/darwin/bin/ticker-term

	tar -cf \
		$(SCRIPTPATH)/build/darwin/ticker-term-$(GIT_TAG)-darwin.tar.gz \
		-C $(SCRIPTPATH)/build/darwin/bin .

build.windows: 
	mkdir -p $(SCRIPTPATH)/build/windows
	GOSUMDB=off GOPROXY=direct CGO_ENABLED=0 GOOS=windows \
		go build -a \
		--ldflags "-s -w \
		-extldflags \"-static\"" \
		-installsuffix cgo \
	    -o $(SCRIPTPATH)/build/windows/ticker-term.exe

debug:
	echo "ARGS $(ARGS)"
	echo "SCRIPTPATH $(SCRIPTPATH)"
	echo "MAKEFILE_PATH $(MAKEFILE_PATH)"
	echo "PRIMARY_REPO $(PRIMARY_REPO)"
	echo "REPO_NAME $(REPO_NAME)"
	echo "REPO_PATH $(REPO_PATH)"
	echo "SSH_URL $(SSH_URL)"
	echo "GO_REPO $(GO_REPO)"
	echo "GO_PKG_PATH $(GO_PKG_PATH)"
	echo "GO_PKG_NAME $(GO_PKG_NAME)"
	echo "BRANCH_NAME $(BRANCH_NAME)"
	echo "COMMIT_SHA $(COMMIT_SHA)"
	echo "GIT_TAG $(GIT_TAG)"

test:
	@echo -e $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

