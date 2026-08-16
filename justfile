set shell := ["zsh", "-eu", "-o", "pipefail", "-c"]

default:
    @just --list

format:
    gofumpt -w .

format-check:
    @unformatted="$(gofumpt -l .)"; if [[ -n "$unformatted" ]]; then print -r -- "$unformatted"; exit 1; fi

test:
    go test ./...

race:
    go test -race ./...

vet:
    go vet ./...

lint:
    staticcheck ./...
    golangci-lint run ./...

security:
    gosec -quiet ./...
    govulncheck ./...

modules:
    go mod verify

build:
    mkdir -p bin
    go build -trimpath -o bin/nb ./cmd/nb

cross-build:
    @tmp="$(mktemp -d "${TMPDIR:-/tmp}/nb-cross.XXXXXX")"; trap 'find "$tmp" -type f -delete; rmdir "$tmp"' EXIT; for target in darwin/arm64 darwin/amd64 linux/arm64 linux/amd64; do os="${target%/*}"; arch="${target#*/}"; echo "building $target"; CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -o "$tmp/nb-$os-$arch" ./cmd/nb; done

licenses:
    ./tools/check-licenses.sh

generated-check:
    ./tools/check-generated.sh

coverage-check:
    ./tools/check-coverage.sh

e2e-selfhosted: build
    ./tools/e2e-selfhosted.sh

diff-check:
    git diff --check

gate: format-check test race vet lint security modules licenses generated-check coverage-check cross-build build diff-check
