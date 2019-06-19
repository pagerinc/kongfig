FROM golang:1.11-stretch@sha256:f7567b33b2319a4e21031f7261ccb1846a0a3870101d2d9e0cdd27cebab29a65

WORKDIR /go/src/github.com/pagerinc/kongfig/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o kongfig

FROM alpine:3.9@sha256:7746df395af22f04212cd25a92c1d6dbc5a06a0ca9579a229ef43008d4d1302a

COPY --from=0 /go/src/github.com/pagerinc/kongfig/kongfig /go/kongfig

RUN apk add --no-cache tini
# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/go/kongfig", "--help"]