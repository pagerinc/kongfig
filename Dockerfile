FROM golang:1.13-stretch@sha256:4e97e0b8479504ee92964c5de6782d88651d953f5df5cfd7a1d45ba48aa90f43

WORKDIR /go/src/github.com/pagerinc/kongfig/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o kongfig

FROM alpine:3.10@sha256:acd3ca9941a85e8ed16515bfc5328e4e2f8c128caa72959a58a127b7801ee01f

COPY --from=0 /go/src/github.com/pagerinc/kongfig/kongfig /go/kongfig

RUN apk add --no-cache tini
# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/go/kongfig", "--help"]