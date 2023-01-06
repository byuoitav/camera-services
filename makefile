NAME := camera-services
OWNER := byuoitav
PKG := github.com/${OWNER}/${NAME}
DOCKER_URL := docker.pkg.github.com
DOCKER_PKG := ${DOCKER_URL}/${OWNER}/${NAME}

# version:
# use the git tag, if this commit
# doesn't have a tag, use the git hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)
TAG := $(shell git rev-parse --short HEAD)
ifneq ($(shell git describe --exact-match --tags HEAD 2> /dev/null),)
	TAG = $(shell git describe --exact-match --tags HEAD)
endif

PRD_TAG_REGEX := "v[0-9]+\.[0-9]+\.[0-9]+"
DEV_TAG_REGEX := "v[0-9]+\.[0-9]+\.[0-9]+-.+"

# go stuff
PKG_LIST := $(shell go list ${PKG}/...)

.PHONY: all deps build test test-cov clean

all: clean build

test:
	@go test -v ${PKG_LIST}

test-cov:
	@go test -coverprofile=coverage.txt -covermode=atomic ${PKG_LIST}

lint:
	@golangci-lint run --tests=false

deps:
	@echo Downloading backend dependencies...
	@go mod download

	@echo Downloading control frontend dependencies...
	@cd cmd/control/web/ && npm install

	@echo Downloading spyglass frontend dependencies...
	@cd cmd/spyglass/web/ && npm install

build: deps
	@mkdir -p dist
	@mkdir -p dist/control
	@mkdir -p dist/spyglass

	@echo
	@echo Building aver for linux-amd64...
	@cd cmd/aver/ && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../../dist/aver-linux-amd64

	@echo
	@echo Building axis for linux-amd64...
	@cd cmd/axis/ && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../../dist/axis-linux-amd64

	#@echo
	#@echo Building slack for linux-amd64...
	#@cd cmd/slack/ && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../../dist/slack-linux-amd64

	@echo
	@echo Building control backend for linux-amd64...
	@cd cmd/control/ && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../../dist/control-linux-amd64

	@echo
	@echo Building control frontend...
	@cd cmd/control/web/ && npm run-script build && ls -la && mv ./dist/web ../../../dist/control/web && rmdir ./dist

	@echo
	@echo Building spyglass backend for linux-amd64...
	@cd cmd/spyglass/ && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../../dist/spyglass-linux-amd64

	@echo
	@echo Building spyglass frontend...
	@cd cmd/spyglass/web/ && npm run-script build && ls -la && mv ./dist/web ../../../dist/spyglass/web && rmdir ./dist

	@echo
	@echo Build output is located in ./dist/.

docker: clean build
ifeq (${COMMIT_HASH}, ${TAG})
	@echo Building dev container with tag ${COMMIT_HASH}

	@echo Building container ${DOCKER_PKG}/aver-dev:${COMMIT_HASH}
	@docker build -f dockerfile --build-arg NAME=aver-linux-amd64 -t ${DOCKER_PKG}/aver-dev:${COMMIT_HASH} dist

	@echo Building container ${DOCKER_PKG}/axis-dev:${COMMIT_HASH}
	@docker build -f dockerfile --build-arg NAME=axis-linux-amd64 -t ${DOCKER_PKG}/axis-dev:${COMMIT_HASH} dist

	#@echo Building container ${DOCKER_PKG}/camera-slack-dev:${COMMIT_HASH}
	#@docker build -f dockerfile --build-arg NAME=slack-linux-amd64 -t ${DOCKER_PKG}/camera-slack-dev:${COMMIT_HASH} dist

	@echo Building container ${DOCKER_PKG}/control-dev:${COMMIT_HASH}
	@docker build -f dockerfile-control --build-arg NAME=control-linux-amd64 -t ${DOCKER_PKG}/control-dev:${COMMIT_HASH} dist

	@echo Building container ${DOCKER_PKG}/spyglass-dev:${COMMIT_HASH}
	@docker build -f dockerfile-spyglass --build-arg NAME=spyglass-linux-amd64 -t ${DOCKER_PKG}/camera-spyglass-dev:${COMMIT_HASH} dist
else ifneq ($(shell echo ${TAG} | grep -x -E ${DEV_TAG_REGEX}),)
	@echo Building dev container with tag ${TAG}

	@echo Building container ${DOCKER_PKG}/aver-dev:${TAG}
	@docker build -f dockerfile --build-arg NAME=aver-linux-amd64 -t ${DOCKER_PKG}/aver-dev:${TAG} dist

	@echo Building container ${DOCKER_PKG}/axis-dev:${TAG}
	@docker build -f dockerfile --build-arg NAME=axis-linux-amd64 -t ${DOCKER_PKG}/axis-dev:${TAG} dist

	#@echo Building container ${DOCKER_PKG}/camera-slack-dev:${TAG}
	#@docker build -f dockerfile --build-arg NAME=slack-linux-amd64 -t ${DOCKER_PKG}/camera-slack-dev:${TAG} dist

	@echo Building container ${DOCKER_PKG}/control-dev:${TAG}
	@docker build -f dockerfile-control --build-arg NAME=control-linux-amd64 -t ${DOCKER_PKG}/control-dev:${TAG} dist

	@echo Building container ${DOCKER_PKG}/spyglass-dev:${TAG}
	@docker build -f dockerfile-spyglass --build-arg NAME=spyglass-linux-amd64 -t ${DOCKER_PKG}/camera-spyglass-dev:${TAG} dist
else ifneq ($(shell echo ${TAG} | grep -x -E ${PRD_TAG_REGEX}),)
	@echo Building prd container with tag ${TAG}

	@echo Building container ${DOCKER_PKG}/aver:${TAG}
	@docker build -f dockerfile --build-arg NAME=aver-linux-amd64 -t ${DOCKER_PKG}/aver:${TAG} dist

	@echo Building container ${DOCKER_PKG}/axis:${TAG}
	@docker build -f dockerfile --build-arg NAME=axis-linux-amd64 -t ${DOCKER_PKG}/axis:${TAG} dist

	#@echo Building container ${DOCKER_PKG}/camera-slack:${TAG}
	#@docker build -f dockerfile --build-arg NAME=slack-linux-amd64 -t ${DOCKER_PKG}/camera-slack:${TAG} dist

	@echo Building container ${DOCKER_PKG}/control:${TAG}
	@docker build -f dockerfile-control --build-arg NAME=control-linux-amd64 -t ${DOCKER_PKG}/control:${TAG} dist

	@echo Building container ${DOCKER_PKG}/spyglass:${TAG}
	@docker build -f dockerfile-spyglass --build-arg NAME=spyglass-linux-amd64 -t ${DOCKER_PKG}/camera-spyglass:${TAG} dist
endif

deploy: docker
	@echo Logging into Github Package Registry
	@docker login ${DOCKER_URL} -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}

ifeq (${COMMIT_HASH}, ${TAG})
	@echo Pushing dev container with tag ${COMMIT_HASH}

	@echo Pushing container ${DOCKER_PKG}/aver-dev:${COMMIT_HASH}
	@docker push ${DOCKER_PKG}/aver-dev:${COMMIT_HASH}

	@echo Pushing container ${DOCKER_PKG}/axis-dev:${COMMIT_HASH}
	@docker push ${DOCKER_PKG}/axis-dev:${COMMIT_HASH}

	#@echo Pushing container ${DOCKER_PKG}/camera-slack-dev:${COMMIT_HASH}
	#@docker push ${DOCKER_PKG}/camera-slack-dev:${COMMIT_HASH}

	@echo Pushing container ${DOCKER_PKG}/control-dev:${COMMIT_HASH}
	@docker push ${DOCKER_PKG}/control-dev:${COMMIT_HASH}

	@echo Pushing container ${DOCKER_PKG}/camera-spyglass-dev:${COMMIT_HASH}
	@docker push ${DOCKER_PKG}/camera-spyglass-dev:${COMMIT_HASH}
else ifneq ($(shell echo ${TAG} | grep -x -E ${DEV_TAG_REGEX}),)
	@echo Pushing dev container with tag ${TAG}

	@echo Pushing container ${DOCKER_PKG}/aver-dev:${TAG}
	@docker push ${DOCKER_PKG}/aver-dev:${TAG}

	@echo Pushing container ${DOCKER_PKG}/axis-dev:${TAG}
	@docker push ${DOCKER_PKG}/axis-dev:${TAG}

	#@echo Pushing container ${DOCKER_PKG}/camera-slack-dev:${TAG}
	#@docker push ${DOCKER_PKG}/camera-slack-dev:${TAG}

	@echo Pushing container ${DOCKER_PKG}/control-dev:${TAG}
	@docker push ${DOCKER_PKG}/control-dev:${TAG}

	@echo Pushing container ${DOCKER_PKG}/camera-spyglass-dev:${TAG}
	@docker push ${DOCKER_PKG}/camera-spyglass-dev:${TAG}
else ifneq ($(shell echo ${TAG} | grep -x -E ${PRD_TAG_REGEX}),)
	@echo Pushing prd container with tag ${TAG}

	@echo Pushing container ${DOCKER_PKG}/aver:${TAG}
	@docker push ${DOCKER_PKG}/aver:${TAG}

	@echo Pushing container ${DOCKER_PKG}/axis:${TAG}
	@docker push ${DOCKER_PKG}/axis:${TAG}

	#@echo Pushing container ${DOCKER_PKG}/camera-slack:${TAG}
	#@docker push ${DOCKER_PKG}/camera-slack:${TAG}

	@echo Pushing container ${DOCKER_PKG}/control:${TAG}
	@docker push ${DOCKER_PKG}/control:${TAG}

	@echo Pushing container ${DOCKER_PKG}/camera-spyglass:${TAG}
	@docker push ${DOCKER_PKG}/camera-spyglass:${TAG}
endif

clean:
	@go clean
	@rm -rf dist/
	@rm -rf cmd/control/web/dist
	@rm -rf cmd/spyglass/web/dist
