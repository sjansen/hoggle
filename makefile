.PHONY: default dist test test-docker

default: test

dist:
	scripts/build-release-binaries

test:
	go test -tags integration ./cmd/... ./pkg/...
	@echo ========================================
	go vet  ./cmd/... ./pkg/...
	golint -set_exit_status cmd/
	golint -set_exit_status pkg/
	gocyclo -over 15 cmd/ pkg/
	@echo ========================================
	@git grep TODO  cmd/ pkg/ || true
	@git grep FIXME cmd/ pkg/ || true

test-docker:
	docker-compose --version
	docker-compose up --abort-on-container-exit --exit-code-from=go --force-recreate
