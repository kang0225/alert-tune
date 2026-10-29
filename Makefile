SHELL := /bin/sh
.DEFAULT_GOAL := check

.PHONY: check fmt-check vet test build smoke

check: fmt-check vet test smoke

fmt-check:
	@set -eu; \
	files="$$(find . -path './.git' -prune -o \
		-path './.claude/worktrees' -prune -o \
		-type f -name '*.go' -exec gofmt -l {} +)"; \
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

# A capture sentinel preserves trailing newlines so extra blank lines are detected.
smoke: build
	@set -eu; \
	version_output="$$(./bin/alert-tune version && printf '.')"; \
	version_output="$${version_output%.}"; \
	flag_output="$$(./bin/alert-tune --version && printf '.')"; \
	flag_output="$${flag_output%.}"; \
	if [ "$$version_output" != "$$flag_output" ]; then \
		printf 'version and --version must print the same output.\n' >&2; \
		exit 1; \
	fi; \
	if ! printf '%s' "$$version_output" | awk \
		'NR == 1 && /^alert-tune [^[:space:]]+$$/ { valid = 1 } END { exit !(NR == 1 && valid) }'; then \
		printf 'Expected a single line containing alert-tune followed by a nonempty version.\n' >&2; \
		exit 1; \
	fi; \
	printf '%s' "$$version_output"
