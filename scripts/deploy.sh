#!/usr/bin/env bash
set -eu

root_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
drone_ip="192.168.1.1"
bin_name="${1:-godrone}"

echo "--> Compiling arm binary ..."
env \
  GOOS=linux \
  GOARCH=arm \
  CGO_ENABLED=0 \
  go install "github.com/felixge/godrone/src/cmd/${bin_name}"

echo "--> Uploading via ftp ..."
curl -T "${root_dir}/gopath/bin/linux_arm/${bin_name}" "ftp://@${drone_ip}/${bin_name}.next"

echo "--> Starting godrone ..."
"${root_dir}/scripts/start.expect" "${drone_ip}" "${bin_name}"
