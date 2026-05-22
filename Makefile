DEBUG_SOCK := /tmp/gotorepl-debug.sock

.PHONY: start
start: ## Start repl

	go run ./cmd/gotorepl/main.go

.PHONY: log
log: ## Start logger

	@nc -lU ${DEBUG_SOCK}
	@echo gotorepl logger started

.PHONY: clean
clean: ## Clean old logger

	rm -f ${DEBUG_SOCK}

# https://postd.cc/auto-documented-makefile/
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:

	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	 sort | \
	 awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'