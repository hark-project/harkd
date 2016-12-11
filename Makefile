GOLINT = golint

lint:
	gofmt -w src/
	$(GOLINT) -set_exit_status src/...

build:
	gb build ...

test:
	gb test ...
