#!/bin/bash

export GOPATH=/go-project
export GOROOT=/usr/local/go
export GOBIN=/usr/local/go/bin
export PATH=/usr/local/go/bin:$PATH
export PATH=/usr/bin:$PATH

export GO111MODULE=off

SOURCE="${BASH_SOURCE[0]}"

while [ -h "${SOURCE}" ]; do
  SCRIPTDIR="$(cd -P "$(dirname "${SOURCE}")" >/dev/null && pwd)"
  SOURCE="$(readlink "${SOURCE}")"
  [[ ${SOURCE} != /* ]] && SOURCE="${SCRIPTDIR}/${SOURCE}"
done
SCRIPTDIR="$(cd -P "$(dirname "${SOURCE}")" >/dev/null && pwd)"

export GOPATH=${GOPATH}

cd ${GOPATH}

arg1=$1
if [ ! -d bin ]
then
    rm -rf bin
    mkdir bin
fi

pwd

rm -rf bin/* && cd bin && go build $arg1
