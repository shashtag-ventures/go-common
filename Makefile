# Makefile for go-common

.PHONY: release check-clean help

# Default release type (major, minor, patch)
TYPE ?= patch

help:
	@echo "Usage:"
	@echo "  make release [PROJECT_DIR=../path/to/project/src] [TYPE=major|minor|patch]"
	@echo ""
	@echo "Examples:"
	@echo "  make release PROJECT_DIR=../karada/src        (default patch)"
	@echo "  make release PROJECT_DIR=../karada/src TYPE=minor"

check-clean:
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: Working directory is not clean. Commit or stash changes first."; \
		exit 1; \
	fi

release: check-clean
	@LATEST_TAG=$$(git tag --sort=-v:refname | head -n 1); \
	if [ -z "$$LATEST_TAG" ]; then echo "Error: No tags found."; exit 1; fi; \
	if [ "$(TYPE)" = "major" ]; then \
		NEXT_VERSION=$$(echo $$LATEST_TAG | awk -F. '{val=substr($$1,2); printf "v%d.0.0\n", val+1}'); \
	elif [ "$(TYPE)" = "minor" ]; then \
		NEXT_VERSION=$$(echo $$LATEST_TAG | awk -F. '{printf "%s.%d.0\n", $$1, $$2+1}'); \
	else \
		NEXT_VERSION=$$(echo $$LATEST_TAG | awk -F. '{printf "%s.%s.%d\n", $$1, $$2, $$3+1}'); \
	fi; \
	echo "Releasing $$NEXT_VERSION (current: $$LATEST_TAG, type: $(TYPE))..."; \
	git tag -a $$NEXT_VERSION -m "Release $$NEXT_VERSION"; \
	git push origin $$NEXT_VERSION; \
	if [ -n "$(PROJECT_DIR)" ]; then \
		echo "Updating project in $(PROJECT_DIR) to $$NEXT_VERSION..."; \
		cd $(PROJECT_DIR) && GOPROXY=direct go get github.com/shashtag-ventures/go-common@$$NEXT_VERSION && go mod tidy; \
	fi; \
	echo "Done!"
