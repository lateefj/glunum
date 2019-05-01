
SHELL := /bin/bash -x
APP := shylock
VERSION := `cat VERSION`

# Support binary builds
PLATFORMS := linux darwin freebsd

all: build

clean:
	rm -fr build
	echo $(PLATFORMS)
	@- $(foreach PLAT,$(PLATFORMS), \
		mkdir -p build/$(PLAT) \
		)
generate:
	rm -f statpkg_gen.go
	go run cmd/codegen/main.go
build: clean
	for plat in $(PLATFORMS); do \
		echo "Building $$plat ..." ; \
		GOARCH=amd64 GOOS=$$plat go build -ldflags "-s -w" -o build/$$plat/$(APP) cmd/shylock/main.go ; \
		done


test:
	go test ./...
