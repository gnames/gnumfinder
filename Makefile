VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')
FLAG_MODULE = GO111MODULE=on
FLAGS_SHARED = $(FLAG_MODULE) CGO_ENABLED=0 GOARCH=amd64
FLAGS_LD=-ldflags "-X github.com/gnames/gnumfinder.Build=${DATE} \
                  -X github.com/gnames/gnumfinder.Version=${VERSION}"
GOCMD=go
GOINSTALL=$(GOCMD) install $(FLAGS_LD)
GOBUILD=$(GOCMD) build $(FLAGS_LD)
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOGET = $(GOCMD) get

all: install

test: deps install
	$(FLAG_MODULE) go test ./...

tools: deps
	@echo Installing tools from tools.go
	@cat gnumfinder/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %


deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download; \

build:
	cd gnumfinder; \
	$(GOCLEAN); \
	$(GOGENERATE) ./...
	$(FLAGS_SHARED) GOOS=linux $(GOBUILD);

release: dockerhub
	cd gnumfinder; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GOBUILD); \
	tar zcvf /tmp/gnumfinder-${VER}-linux.tar.gz gnumfinder; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=darwin $(GOBUILD); \
	tar zcvf /tmp/gnumfinder-${VER}-mac.tar.gz gnumfinder; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=windows $(GOBUILD); \
	zip -9 /tmp/gnumfinder-${VER}-win-64.zip gnumfinder.exe; \
	$(GOCLEAN);

install:
	cd gnumfinder; \
	$(FLAGS_SHARED) $(GOINSTALL);

docker: build
	docker build -t gnames/gnumfinder:latest -t gnames/gnumfinder:${VERSION} .; \
	cd gnumfinder; \
	$(GOCLEAN);

dockerhub: docker
	docker push gnames/gnumfinder; \
	docker push gnames/gnumfinder:${VERSION}
