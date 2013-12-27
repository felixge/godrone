version = $(shell git describe --tags --dirty)

dev:
	./scripts/build.bash -version $(version) dist/dev
	./dist/dev/deploy

dev-local:
	go build github.com/felixge/godrone/cmd/godrone
	./godrone cmd/godrone/godrone.conf

dist:
	./scripts/dist.bash $(version)

.PHONY: dev dist
