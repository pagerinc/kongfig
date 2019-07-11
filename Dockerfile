FROM golang:1.12-stretch@sha256:00687f5beb8b17b8af7a756943fb1658832615d714c1760562f5ae40492e1107

WORKDIR /go/src/github.com/pagerinc/kongfig/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o kongfig

FROM alpine:3.10@sha256:6a92cd1fcdc8d8cdec60f33dda4db2cb1fcdcacf3410a8e05b3741f44a9b5998

COPY --from=0 /go/src/github.com/pagerinc/kongfig/kongfig /go/kongfig

RUN apk add --no-cache tini
# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/go/kongfig", "--help"]