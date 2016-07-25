PACKAGES = $(shell go list ./... | grep -v /vendor/)

all: build

deps:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...
	go get -u github.com/vektra/mockery/...

gen: gen_assets gen_mocks gen_migration

gen_assets:
	go generate github.com/lgtmco/lgtm/web/static
	go generate github.com/lgtmco/lgtm/web/template

gen_mocks:
	go generate github.com/lgtmco/lgtm/store
	go generate github.com/lgtmco/lgtm/remote

gen_migration:
	go generate github.com/lgtmco/lgtm/store/migration

build:
	@# do not staticly link sqlite. Doesn't work on OS X and breaks DNS resolution on linux
	go build --ldflags '-X github.com/lgtmco/lgtm/version.VersionDev=$(CI_BUILD_NUMBER)' -o lgtm

test:
	@for PKG in $(PACKAGES); do go test -v -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG; done;

test_mysql:
	DATABASE_DRIVER="mysql" DATABASE_DATASOURCE="root@tcp(127.0.0.1:3306)/test?parseTime=true" go test -v -cover github.com/lgtmco/lgtm/store/datastore
