FROM golang:1.9 as builder
WORKDIR /go/src/github.com/mvisonneau/strongbox
COPY . .
RUN
make prereqs ;\
make dep ;\
make build

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/mvisonneau/strongbox/strongbox /
ENTRYPOINT ["/strongbox"]
CMD [""]
