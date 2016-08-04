# Build the drone executable on a x64 Linux host:
#
#     go build --ldflags '-extldflags "-static"' -o lgtm
#
# Build the docker image:
#
#     docker build --rm=true -t lgtm/lgtm .
#
# Push to Heroku
#
# 		heroku container:push web

FROM golang:1.6

EXPOSE 8000

ENV DATABASE_DRIVER=sqlite3
ENV DATABASE_DATASOURCE=/var/lib/lgtm/lgtm.sqlite
ENV GO15VENDOREXPERIMENT=1

COPY . /go/src/github.com/lgtmco/lgtm
WORKDIR /go/src/github.com/lgtmco/lgtm

RUN make deps
RUN make gen
RUN make build

CMD ["/go/src/github.com/lgtmco/lgtm/lgtm"]
