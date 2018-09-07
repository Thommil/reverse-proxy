#!/bin/sh

if [ -z "$GOPATH" ] ; then
    export GOPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    export PATH=$PATH:$GOPATH/bin
fi

# Deps
go get github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/protoc-gen-micro
go get github.com/favadi/protoc-go-inject-tag

#Build proto files
protoc --plugin=protoc-gen-go=$GOPATH/bin/protoc-gen-go --plugin=protoc-gen-micro=$GOPATH/bin/protoc-gen-micro --proto_path=$GOPATH/src:. --micro_out=. --go_out=. proto/user.proto

#Inject tags in generated files
protoc-go-inject-tag -input=$GOPATH/src/gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto/user.pb.go

#build
go install gitlab.forge.orange-labs.fr/mahali/services/ms-user