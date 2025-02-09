
@PHONY: build
build:
	@echo "Building..."
	cd scraper && \
	go build -o scraper main.go && \
	chmod u+x scraper
	@echo "Done."

@PHONY: generate
generate:
	@echo "Generating..."
	cd scraper && \
	./scraper -u carlos-algms
