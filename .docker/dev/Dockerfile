FROM golang:1.9
WORKDIR /go/src/github.com/mvisonneau/strongbox
COPY . .
RUN \
make prereqs ;\
make dep ;\
make install
ENTRYPOINT ["strongbox"]
CMD [""]
