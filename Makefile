PACKAGES = $(shell go list ./... | grep -v /vendor/)

all: build

deps:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...
	go get -u github.com/vektra/mockery/...

gen:
	go generate github.com/lgtmco/lgtm/web/static
	go generate github.com/lgtmco/lgtm/web/template
	go generate github.com/lgtmco/lgtm/notifier
	go generate github.com/lgtmco/lgtm/remote
	go generate github.com/lgtmco/lgtm/store/migration
	go generate github.com/lgtmco/lgtm/store

build:
	go build --ldflags '-extldflags "-static" -X github.com/lgtmco/lgtm/version.VersionDev=$(CI_BUILD_NUMBER)' -o lgtm

test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG; done;

test_mysql:
	DATABASE_DRIVER="mysql" DATABASE_DATASOURCE="root@tcp(127.0.0.1:3306)/test?parseTime=true" go test -v -cover github.com/lgtmco/lgtm/store/datastore
