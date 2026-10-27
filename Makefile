SHELL := /bin/sh
.DEFAULT_GOAL := check

.PHONY: check fmt-check vet test build smoke

check: fmt-check vet test smoke

fmt-check:
	@set -eu; \
	files="$$(gofmt -l .)"; \
	if [ -n "$$files" ]; then \
		printf 'Run gofmt on these files:\n%s\n' "$$files" >&2; \
		exit 1; \
	fi

vet:
	go vet ./...

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/alert-tune ./cmd/alert-tune

smoke: build
	@set -eu; \
	version_output="$$(./bin/alert-tune version)"; \
	flag_output="$$(./bin/alert-tune --version)"; \
	if [ "$$version_output" != "$$flag_output" ]; then \
		printf 'version and --version must print the same output.\n' >&2; \
		exit 1; \
	fi; \
	if ! printf '%s\n' "$$version_output" | grep -Eq '^alert-tune [^[:space:]]+$$'; then \
		printf 'Expected alert-tune followed by a nonempty version.\n' >&2; \
		exit 1; \
	fi; \
	printf '%s\n' "$$version_output"
