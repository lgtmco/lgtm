# Build the drone executable on a x64 Linux host:
#
#     go build --ldflags '-extldflags "-static"' -o lgtm
#
# Build the docker image:
#
#     docker build --rm=true -t lgtm/lgtm .

FROM centurylink/ca-certs
EXPOSE 8000

ENV DATABASE_DRIVER=sqlite3
ENV DATABASE_DATASOURCE=/var/lib/lgtm/lgtm.sqlite
ENV GODEBUG=netdns=go

ADD lgtm /lgtm
ENTRYPOINT ["/lgtm"]
