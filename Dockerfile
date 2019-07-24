FROM golang:1.12-stretch@sha256:88e108f3f5410e9e184bd61493306ab0be7e494fb5e4c42c33fc18284ef3a222

WORKDIR /go/src/github.com/pagerinc/kongfig/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o kongfig

FROM alpine:3.10@sha256:6a92cd1fcdc8d8cdec60f33dda4db2cb1fcdcacf3410a8e05b3741f44a9b5998

COPY --from=0 /go/src/github.com/pagerinc/kongfig/kongfig /go/kongfig

RUN apk add --no-cache tini
# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/go/kongfig", "--help"]