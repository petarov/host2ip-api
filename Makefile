# go compile stuff
BIN_NAME_PREFIX	= host2ip
BULD_FLAGS		=
LDLFAGS			= "-s -w"
# docker stuff
DOCKER		= docker
CONTAINER	= host2ip
IMG_NAME	= github.com/petarov/host2ip-api
IMG_TAG		= latest


.PHONY: all

build: main.main
all: clean build dist

%.main: main.go
	GOOS=linux GOARCH=amd64 go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(BIN_NAME_PREFIX)_linux_amd64
	GOOS=windows GOARCH=amd64 go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(BIN_NAME_PREFIX)_windows_amd64.exe

.PHONY: dist
dist:
	test -d dist || mkdir -p dist/
	mv $(BIN_NAME_PREFIX)_* dist/

.PHONY: clean
clean:
	@rm -f $(BIN_NAME_PREFIX)_*
	@test -d dist && @rm -f dist/$(BIN_NAME_PREFIX)~
	@test -d dist && @rmdir dist

.PHONY: build-docker
build-docker:
	$(DOCKER) build . -t ${IMG_NAME}:${IMG_TAG}

.PHONY: run-docker
run-docker:
	$(DOCKER) run -p 7029:7029/tcp --name ${CONTAINER} ${IMG_NAME}:${IMG_TAG}