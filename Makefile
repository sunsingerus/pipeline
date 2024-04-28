# Absolute path to project root
PROJECT_ROOT=$(shell dirname $(abspath $(lastword ${MAKEFILE_LIST})))
# Absolute path to project's 'tmp' folder
PROJECT_TMP=${PROJECT_ROOT}/tmp
# Absolute path to current makefile
MAKEFILE=$(abspath $(lastword ${MAKEFILE_LIST}))
# Absolute path to 'main' module
MAIN_MODULE_PATH=${PROJECT_ROOT}/cmd/pipeline
# Executable name
EXE_PATH=${PROJECT_ROOT}/bin/pipeline

# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

## help: show available targets with short descriptions
.PHONY: help
help:
	@echo "Launching target: ${@}"
	@sed -n 's/^##//p' ${MAKEFILE} | column -t -s ':' |  sed -e 's/^/ /'

## release: show project release
.PHONY: release
release:
	@echo "Launching target: ${@}"
	@cat release

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty-files
no-dirty-files:
	git diff --exit-code

## format: format code
.PHONY: format
format:
	@echo "Launching target: ${@}"
	go fmt ${PROJECT_ROOT}/...

## tidy: format code and tidy go.mod file
.PHONY: tidy
tidy: format
	@echo "Launching target: ${@}"
	go mod tidy -v

## vendor: vendor dependencies
.PHONY: vendor
vendor: tidy
	@echo "Launching target: ${@}"
	rm -rf ${PROJECT_ROOT}/vendor
	go mod vendor

## audit: run quality control checks
.PHONY: audit
audit:
	@echo "Launching target: ${@}"
	go mod verify
	go vet ${PROJECT_ROOT}/...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ${PROJECT_ROOT}/...
	#go run golang.org/x/vuln/cmd/govulncheck@latest ${PROJECT_ROOT}/...
	go test -race -buildvcs -vet=off ${PROJECT_ROOT}/...

## test: run all tests
.PHONY: test
test:
	@echo "Launching target: ${@}"
	go test -v -race -buildvcs ${PROJECT_ROOT}/...

COVERAGE_FILE=${PROJECT_TMP}/coverage.out
## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	@echo "Launching target: ${@}"
	@[ -f ${COVERAGE_FILE} ] && rm -rf ${COVERAGE_FILE}
	go test -v -race -buildvcs -coverprofile=${COVERAGE_FILE} ${PROJECT_ROOT}/...
	@[ -f ${COVERAGE_FILE} ] && go tool cover -html=${COVERAGE_FILE}

## build: build the application
.PHONY: build
build:
	@echo "Launching target: ${@}"
	go build -o=${EXE_PATH} ${MAIN_MODULE_PATH}

## run: run the  application
.PHONY: run
run: build
	@echo "Launching target: ${@}"
	@echo "Run app:"
	${EXE_PATH} $(RUN_ARGS)

## git/push: push changes to the remote Git repository
.PHONY: git/push
git/push: vendor audit no-dirty-files
	@echo "Launching target: ${@}"
	git push

## deploy: deploy the application
.PHONY: deploy
deploy: confirm tidy audit no-dirty-files
	@echo "Launching target: ${@}"
