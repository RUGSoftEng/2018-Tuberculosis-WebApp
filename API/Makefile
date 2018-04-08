GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
EXEC_NAME=API

all: format check build
run: build
	./$(EXEC_NAME)

build:
	$(GOBUILD) -o $(EXEC_NAME)

#########################
# Go Tools for formatting
format:	fmt lint

# -w: writes changes to the file
# -l: list changed files
fmt:
	gofmt -l -w *.go

# -set_exit_status: returns 1 if any improvements could be found
lint:
	golint -set_exit_status *.go

#############################
# Go Tools for error checking
check: vet sql

vet:
	go vet *.go

# Looks for SQL Injection weaknesses
# -q: No output if nothing fails
sql:
	safesql -q .

#####################################################
# Dependencies for running the makefile
# This does not get the dependencies for the api code
deps:
	go get -u golang.org/x/lint/golint
	go get -u github.com/stripe/safesql

clean:
	rm -f $(EXEC_NAME)
	rm -f *~ # Emacs backup files
