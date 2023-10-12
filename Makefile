# Version: 1.0.0

# Go parameters
GOCMD=go
GOBUILD=CGO_ENABLED=1 CGO_CFLAGS="-g -O2 -Wno-return-local-addr" $(GOCMD) build -v
GOTESTS=CGO_ENABLED=1 CGO_CFLAGS="-g -O2 -Wno-return-local-addr" $(GOCMD) test
GOCLEAN=$(GOCMD) clean

VERSION=1.3.0

DIRECTORY_BIN=bin
DIRECTORY_TMP=./resources/templates

all: create-templates copy test build

create-templates:
	# Create e-mail templates
	cd templates && npm install --save-dev && npm run build

copy:
	# Create templates directory
	if [ ! -d "$(DIRECTORY_TMP)" ]; then mkdir -p $(DIRECTORY_TMP); fi

	# Copy templates to templates directory
	cp -r ./templates/dist/* $(DIRECTORY_TMP)

test:
	# Run tests
	$(GOTESTS) ./...

build:
	# Create bin directory
	if [ ! -d "./$(DIRECTORY_BIN)" ]; then mkdir $(DIRECTORY_BIN); fi

	# Native build
	$(GOBUILD) $(BUILDFLAGS) -o $(DIRECTORY_BIN)/gobarista pkg/cmd/gobarista/*.go
	
	# Linux AMD64 architecture build
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILDFLAGS) -o $(DIRECTORY_BIN)/gobarista-$(VERSION).linux-amd64 pkg/cmd/gobarista/*.go

clean:
	rm -rf $(DIRECTORY_BIN)
