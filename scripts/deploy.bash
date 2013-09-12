#!/usr/bin/env bash
set -eu

root_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
pkg_path="github.com/felixge/godrone/cmd"
bin_name="godrone"
drone_ip="${1:-192.168.1.1}"

echo "--> Fetching dependencies ..."
go get "${pkg_path}"

echo "--> Compiling arm binary ..."
env \
  GOOS=linux \
  GOARCH=arm \
  CGO_ENABLED=0 \
  go build -o "${bin_name}" "${pkg_path}"

echo "--> Uploading via ftp ..."
curl -T "${bin_name}" "ftp://@${drone_ip}/${bin_name}.next"

echo "--> Starting godrone ..."
"${root_dir}/scripts/start.expect" "${drone_ip}" "${bin_name}"
