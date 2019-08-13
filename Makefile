USERNAME="reaandrew"
PROJECT="skeleton-go-system"
GITHUB_TOKEN=$$GITHUB_TOKEN
VERSION=`cat VERSION`
BUILD_TIME=`date +%FT%T%z`
COMMIT_HASH=`git rev-parse HEAD`
DIST_NAME_CONVENTION="dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
SOURCES += VERSION
# This is how we want to name the binary output
BINARY=${PROJECT}

# These are the values we want to pass for Version and BuildTime

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.CommitHash=${COMMIT_HASH} -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): deps $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} 

.PHONY: deps 
instal-deps:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -t -v -d ./...

.PHONY: deploy-deps
install-deploy-deps:
	go get -u github.com/mitchellh/gox
	go get -u github.com/tcnksm/ghr
	go get -u github.com/mattn/goveralls
	GOOS=windows go get -u github.com/spf13/cobra

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -coverprofile=coverage.out ./...

.PHONY: coverage-report
coverage-report: 
	go tool cover -html=c.out -o coverage.html

.PHONY: cross-platform-compile
cross-platform-compile: deploy-deps
	gox -output ${DIST_NAME_CONVENTION} ${LDFLAGS}

.PHONY: upload-release
upload-release:
	ghr -username ${USERNAME} -token ${GITHUB_TOKEN} dist/

