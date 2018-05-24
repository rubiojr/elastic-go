EXEC = esg
PKG  = github.com/rubiojr/esg

all: check test build

build: clean
	go build -o ${EXEC} ${PKG}

install: clean
	go install ${PKG}

test:
	go test -v ${PKG}/...

cover:
	go test -cover ${PKG}/...

check:
	go vet ${PKG}/...
	golint

deps:
	go get -u github.com/gilliek/go-xterm256/xterm256
	go get -u github.com/hokaccha/go-prettyjson
	go get -u github.com/rivo/tview
	go get -u github.com/urfave/cli

deps-dev: deps
	go get -u -v github.com/golang/lint/golint

clean:
	rm -f ${EXEC}

.PHONY: build install test cover check deps deps-dev clean

