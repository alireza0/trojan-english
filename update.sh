#!/bin/bash

#Get the absolute path where the current script is located
SHELL_PATH=$(cd `dirname $0`; pwd)

cd $SHELL_PATH

mkdir -p web/templates

touch web/templates/test

go get -u

rm -rf web/templates