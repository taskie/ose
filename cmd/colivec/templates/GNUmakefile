.PHONY: build
build:
	go build -v -ldflags "-s -w"
{{ range .CmdNames }}	cd cmd/{{ . }} && go build -v -ldflags "-s -w"
{{ end }}

.PHONY: install
install:
	go install -v -ldflags "-s -w"
{{ range .CmdNames }}	cd cmd/{{ . }} && go install -v -ldflags "-s -w"
{{ end }}

.PHONY: generate
	go generate

.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	mkdir -p test/coverage
	go test -coverprofile=test/coverage/cover.out ./...
	go tool cover -html=test/coverage/cover.out -o test/coverage/cover.html

.PHONY: fmt
fmt:
	find . -name '*.go' | grep -v ./vendor/ | xargs gofmt -w

.PHONY: upgrade
upgrade:
	go get -u
{{ range .CmdNames }}	cd cmd/{{ . }} && go get -u
{{ end }}
	$(MAKE) mod-tidy

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: pre-commit
pre-commit:
	$(MAKE) mod-tidy
	$(MAKE) fmt
	$(MAKE) build
	$(MAKE) test

.PHONY: install-pre-commit
install-pre-commit:
	echo 'make pre-commit' >.git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
