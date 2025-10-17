#!/bin/bash
set -e
workDir=$(cd $(dirname $0) && pwd)
cd $workDir
make build
mkdir output
app_name=avatar
cp -v $app_name prod.yaml test.yaml start.sh output 
