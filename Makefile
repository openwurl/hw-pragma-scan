#########################
###      DEFS         ###
#########################

# Don't ask, for to understand it is to look 
# into the void and know the void is not only 
#looking back but also reading your emails.
SHELL=/bin/bash -e -o pipefail
ARCH?=amd64
OS?=darwin

#########################
###     Targets       ###
#########################

.PHONY: build-osx
.DEFAULT_GOAL := build-osx

dep:
	@go mod tidy

build: clean
	GOARCH=$(ARCH) GOOS=$(OS) go build -o build/hw-pragma-$(OS)

clean:
	@rm -rf build/
	@go clean