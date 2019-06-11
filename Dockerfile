##
# BUILD CONTAINER
##

FROM golang:1.12 as builder

WORKDIR /build

COPY Makefile .
RUN \
make setup

COPY . .
RUN \
make build-docker

##
# RELEASE CONTAINER
##

FROM busybox:1.31-glibc

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/strongbox /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/strongbox"]
CMD [""]
