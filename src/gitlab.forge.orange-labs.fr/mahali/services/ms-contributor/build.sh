#!/bin/sh

if [ -z "$GOPATH" ]; then
    export GOPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    export PATH=$PATH:$GOPATH/bin
fi    

#build
go install gitlab.forge.orange-labs.fr/mahali/services/ms-contributor