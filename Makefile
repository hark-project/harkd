GOLINT = golint
GB = gb
GOFMT = gofmt

SRCDIR = src

TESTFLAGS = -race

all: lint test build

lint:
	$(GOFMT) -w $(SRCDIR)
	$(GOLINT) -set_exit_status $(SRCDIR)/...

build:
	$(GB) build

test:
	$(GB) test $(TESTFLAGS)
