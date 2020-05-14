#!/usr/bin/env bash

DIR=`dirname ${0}`
. $DIR/common.sh

$BASE_DIR/gow test -v -cover $BASE_DIR/...