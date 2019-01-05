test:
	@echo 'Running test...'
	@go test -v -cover ./... || echo 'test failed :('

	@echo 'Clearing formal test data...'
	@rm -r ~/.gitdoc &> /dev/null || true
