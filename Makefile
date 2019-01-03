test:
	@rm -r ~/.gitdoc || true
	@mkdir ~/.gitdoc || true
	@go test -v -cover ./... || echo 'test failed :('
