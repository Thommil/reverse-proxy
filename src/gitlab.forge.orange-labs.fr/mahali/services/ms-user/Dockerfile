FROM golang:latest

WORKDIR /go
ENV GOPATH /go

ADD src /go/src
RUN go get -v gitlab.forge.orange-labs.fr/mahali/services/ms-user

# To switch between HTTP and gRPC plugins, set env. :
#   MICRO_CLIENT=grpc | http
#   MICRO_SERVER=grpc | http
#
# !!! Be careful when mixing types !!!

CMD ["ms-user"]