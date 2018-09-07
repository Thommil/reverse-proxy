FROM golang:latest

WORKDIR /go
ENV GOPATH /go

ADD src /go/src
RUN go install -v gitlab.forge.orange-labs.fr/mahali/services/ms-contributor

# To switch between HTTP and gRPC plugins, set env. :
#   MICRO_CLIENT=grpc | http

CMD ["ms-contributor"]