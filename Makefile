version = $(shell git describe --tags --dirty)
host_os = $(shell go env GOOS)
host_arch = $(shell go env GOARCH)

dev:
	./scripts/build.bash -version $(version) dist/dev
	./dist/dev/deploy

dev-local:
	env \
		BUILD_OS="$(host_os)" \
		BUILD_ARCH="$(host_arch)" \
		./scripts/build.bash -version $(version) dist/local-dev
	cd dist/local-dev && ./godrone

dist:
	./scripts/dist.bash $(version)

.PHONY: dev dist
