##
# BUILD CONTAINER
##

FROM golang:1.12 as builder

WORKDIR /go/src/github.com/mvisonneau/strongbox

COPY Makefile .
RUN \
make setup

COPY . .
RUN \
make build

##
# RELEASE CONTAINER
##

FROM busybox:1.28-glibc

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/mvisonneau/strongbox/strongbox /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/strongbox"]
CMD [""]
