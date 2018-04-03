GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
EXEC_NAME=API

all: build check
run: build
	./$(EXEC_NAME)

build:
	$(GOBUILD) -o $(EXEC_NAME)

check: fmt lint vet

# w: writes changes to the file
# l: list changed files
fmt:
	gofmt -l -w *.go

# -set_exit_status: returns 1 if any improvements could be found
lint:
	golint -set_exit_status *.go

vet:
	go vet *.go

clean:
	rm -f $(EXEC_NAME)
	rm -f *~ # Emacs backup files
