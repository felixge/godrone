version = $(shell git describe --tags --dirty)

dev:
	./scripts/build.bash -version $(version) dist/dev
	./dist/dev/deploy

dist:
	./scripts/dist.bash $(version)

.PHONY: dev dist
