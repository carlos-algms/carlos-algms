
build:
	@echo "Building..."
	cd scraper && \
	go build -o scraper main.go && \
	chmod u+x scraper
	@echo "Done."
