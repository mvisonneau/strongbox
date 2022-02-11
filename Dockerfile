##
# BUILD CONTAINER
##

FROM alpine:3.15.0 as certs

RUN \
apk add --no-cache ca-certificates

##
# RELEASE CONTAINER
##

FROM busybox:1.35-glibc

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY strongbox /usr/local/bin/

# Run as nobody user
USER 65534

ENTRYPOINT ["/usr/local/bin/strongbox"]
CMD [""]
