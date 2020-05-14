#!/usr/bin/env bash

DIR=`dirname ${0}`
. $DIR/common.sh


$BASE_DIR/gow build -o $BASE_DIR/bin/go-boot-web$extension $BASE_DIR/cmd/go-boot-web/main.go