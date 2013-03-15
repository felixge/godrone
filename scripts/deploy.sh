#!/usr/bin/env bash
set -eu

root_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
drone_ip="192.168.1.1"

# Environment variables we need to cross compile for the drone
#export GOOS := linux
#export GOARCH := arm
#export CGO_ENABLED=0

#define run
#curl -T $1 ftp://@$(DRONE_IP)/upload
#(\
#echo spawn telnet $(DRONE_IP);\
#echo expect -re .*#;\
#echo send \"cd /data/video\\\r\";\
#echo expect -re .*#;\
#echo send \"killall $1\\\r\";\
#echo expect -re .*#;\
#echo send \"rm $1\\\r\";\
#echo expect -re .*#;\
#echo send \"mv upload $1\\\r\";\
#echo expect -re .*#;\
#echo send \"chmod +x $1\\\r\";\
#echo expect -re .*#;\
#echo send \"./$1\\\r\";\
#echo set timeout -1;\
#echo expect -re .*#;\
#) | expect
#endef

echo "--> Compiling arm binary ..."
env \
  GOOS=linux \
  GOARCH=arm \
  CGO_ENABLED=0 \
  go install github.com/felixge/godrone/src/cmd/godrone

echo "--> Uploading via ftp ..."
curl -T "${root_dir}/gopath/bin/linux_arm/godrone" "ftp://@${drone_ip}/upload"
