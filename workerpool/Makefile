test:
	@go test ./...

cover_html:
	@mkdir -p _out
	@go test -coverprofile _out/cover.out ./...
	@go tool cover -html _out/cover.out -o ./_out/coverage.html
	@open ./_out/coverage.html&
