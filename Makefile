DEP_EXISTS := $(shell command -v dep 2> /dev/null)

.PHONY: all
all: build

.PHONY: build
build: ensure_deps compile

.PHONY: compile
compile: 
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

.PHONY: ensure_deps
ensure_deps: install_dep
	dep ensure -v

.PHONY: install_dep
install_dep:
ifndef DEP_EXISTS
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: docker
docker:
	docker build --tag roaanv/k8sti .
