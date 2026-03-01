BUILD_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT        := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE          := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS       := -X main.version=$(BUILD_VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

BINARY  := recon
BIN_DIR := bin

.PHONY: build test lint install clean release

build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) ./cmd/recon

test:
	go test -race ./...

lint:
	golangci-lint run ./...

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/recon

clean:
	rm -rf $(BIN_DIR)

# make release VERSION=1.2.3
# Updates CHANGELOG, commits, and tags. Then run: git push origin main vVERSION
release:
ifndef VERSION
	$(error VERSION is not set. Usage: make release VERSION=1.2.3)
endif
	@echo "Preparing release v$(VERSION)..."
	@RELEASE_DATE=$$(date -u +%Y-%m-%d); \
	perl -i -0pe "s/## \[Unreleased\]/## [Unreleased]\n\n## [$(VERSION)] - $$RELEASE_DATE/" CHANGELOG.md
	@perl -i -pe \
		's|\[Unreleased\]: (https://.+)/compare/v.+\.\.\.HEAD|[Unreleased]: $$1/compare/v$(VERSION)...HEAD\n[$(VERSION)]: $$1/releases/tag/v$(VERSION)|' \
		CHANGELOG.md
	git add CHANGELOG.md
	git commit -m "chore: release v$(VERSION)"
	git tag -d v$(VERSION) 2>/dev/null || true
	git tag v$(VERSION)
	@echo ""
	@echo "Release v$(VERSION) prepared. To publish:"
	@echo "  git push origin master v$(VERSION)"

.DEFAULT_GOAL := build
