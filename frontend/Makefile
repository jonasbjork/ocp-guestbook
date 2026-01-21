SHELL := /usr/bin/env bash
HTMLHINT := $(shell which htmlhint)

.PHONY: check lint

check:
	@if [ -z "${HTMLHINT}" ]; then \
	    echo "htmlhint saknas"; \
	    exit 1; \
	fi

lint: check
	htmlhint index.html


