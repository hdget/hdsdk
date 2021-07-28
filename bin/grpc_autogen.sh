#!/usr/bin/env bash

BASE_DIR=$(dirname $(dirname $(readlink -f "$0")))
RPC_DIR=$BASE_DIR/rpc

GOPATH=`go env GOPATH`
PATH=$PATH:$BASE_DIR/bin:$GOPATH/bin

SERVERS=("treeserver")
for s in $SERVERS; do
  DIR=$BASE_DIR/$s/service
  for i in `find $DIR -type d`; do
    cd $i
    go generate -v
    cd $BASE_DIR
  done
done

