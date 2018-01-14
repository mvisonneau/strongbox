##
# BUILD CONTAINER
##

FROM golang:1.9.2 as builder

WORKDIR /go/src/github.com/mvisonneau/strongbox

COPY Makefile .
RUN \
make prereqs

COPY . .
RUN \
make deps ;\
make build

##
# RELEASE CONTAINER
##

FROM scratch

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/mvisonneau/strongbox/strongbox /

ENTRYPOINT ["/strongbox"]
CMD [""]
