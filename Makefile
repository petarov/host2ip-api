EXEC_NAME_PREFIX=host2ip
BULD_FLAGS=
LDLFAGS="-s -w"

.PHONY: all

build: main.main
all: clean build dist

%.main: main.go
	GOOS=linux GOARCH=amd64 go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(EXEC_NAME_PREFIX)_linux_amd64
	GOOS=windows GOARCH=amd64 go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(EXEC_NAME_PREFIX)_windows_amd64.exe

dist:
	test -d dist || mkdir -p dist/
	mv $(EXEC_NAME_PREFIX)_* dist/

clean:
	@rm -f $(EXEC_NAME_PREFIX)_*
	@test -d dist && @rm -f dist/$(EXEC_NAME_PREFIX)~
	@test -d dist && @rmdir dist