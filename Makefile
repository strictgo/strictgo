.PHONY: build install install-tools test vet lint fmt tidy check-tidy git-clean check ci clean

build:
	go build -o bin/strictgo ./cmd/strictgo

install:
	go install ./cmd/strictgo

install-tools:
	$(MAKE) -C tools install

test:
	go test ./...

vet:
	go vet ./...

lint: install-tools vet build
	./bin/staticcheck ./...
	./bin/errcheck ./...
	./bin/strictgo ./...

fmt: install-tools
	go fix ./...
	./bin/gofumpt -w .

tidy:
	go mod tidy
	cd tools && go mod tidy

check-tidy:
	go mod tidy -diff
	cd tools && go mod tidy -diff

git-clean:
	@status="$$(git status --porcelain=v1 --untracked-files=all)" || exit 1; \
	if [ -n "$$status" ]; then \
		printf 'Working tree is dirty:\n%s\n' "$$status" >&2; \
		exit 1; \
	fi

check: check-tidy test lint

ci: check
	$(MAKE) git-clean

clean:
	rm -rf bin
