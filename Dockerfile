##
# BUILD CONTAINER
##

FROM goreleaser/goreleaser:v0.147.2 as builder

WORKDIR /build

COPY Makefile .
RUN \
apk add --no-cache make ;\
make setup

COPY . .
RUN \
make build-linux-amd64

##
# RELEASE CONTAINER
##

FROM busybox:1.32.0-glibc

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/dist/strongbox_linux_amd64/strongbox /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/strongbox"]
CMD [""]
