BUILDDIR ?= bin
PWD ?= `pwd`
ADDFLAGS ?=
BUILDFLAGS ?= $(ADDFLAGS) -ldflags '-w -s'
CGOFLAG ?= CGO_ENABLED=1
GOPKGDIR ?= `go env GOPATH`/pkg/`go env GOHOSTOS`_`go env GOARCH`/github.com/marve39/tvloka/c1*

ARCH ?= $(shell go env GOARCH)
OS ?= $(shell go env GOOS)

ifeq ("$(OS)","windows")
BUILDFLAGS = $(ADDFLAGS) -ldflags '-w -s' -buildmode=exe
CGOFLAG = CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++
endif

ifeq ("$(OS)","linux")
# ARM builds need to specify the correct C compiler
ifeq ("$(ARCH)","arm")
CGOFLAG = CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc
endif
# ARM64 builds need to specify the correct C compiler
ifeq ("$(ARCH)","arm64")
CGOFLAG = CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc
endif
endif

.PHONY: run-ci
run-ci:
	@echo "---> Running CI pipeline"
	${MAKE} ci-lint 
	${MAKE} ci-vet 
	${MAKE} ci-test 

.PHONY: ci-lint
ci-lint:
	@echo "---> Running lint check"
	golint ./...

.PHONY: ci-vet
ci-vet:
	@echo "---> Running vet check"
	go vet ./...

.PHONY: ci-test
ci-test:
	@echo "---> Running test"
	go test -v ./... -coverprofile=coverage.out -json | tee test-result.json

.PHONY: ci-test-unit
ci-test-unit:
	@echo "---> Running unit test"
	go test -v ./... -run ^TestUnit$

.PHONY: ci-test-integration
ci-test-integration:
	@echo "---> Running integration test"
	go test -v ./... -run ^TestIntegration$

.PHONY: build
build:
	@echo "---> Building binaries."
	GOOS=$(OS) GOARCH=$(ARCH) $(CGOFLAG) go build -o $(BUILDDIR)/c1-bin $(BUILDFLAGS) ./main.go

.PHONY: clean
clean:
	@echo "---> Cleaning build artifacts."
	go mod tidy
	rm -rf $(BUILDDIR)
	rm -rf coverage.*
	rm -rf test-result.*
	go clean -cache
	rm -rf $(GOPKGDIR)

.PHONY: docker-build
docker-build:
	@echo "---> Building Docker Image."
	DOCKER_BUILDKIT=1 docker build --progress plain --tag c1-bin:latest --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} --file Dockerfile .
	mkdir -p tmp_trivy
	docker run --rm -v ${PWD}/tmp_trivy:/root/.cache/ -v /var/run/docker.sock:/var/run/docker.sock \
    aquasec/trivy --exit-code 1 --severity CRITICAL c1-bin:latest
	rm -rf tmp_trivy

.PHONY: docker-run
docker-run:
	@echo "---> Running Docker Image."
	docker run --rm -e GITHUB_TOKEN=${GITHUB_TOKEN} -it c1-bin 