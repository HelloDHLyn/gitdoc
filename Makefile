test:
	@echo 'Clearing formal test data...'
	@rm -r ~/.gitdoc &> /dev/null || true
	@mkdir ~/.gitdoc &> /dev/null || true

	@echo 'Running test...'
	@go test -v -cover ./... || echo 'test failed :('
