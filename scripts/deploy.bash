#!/usr/bin/env bash
set -eu

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
root_dir="$( cd "${dir}" && cd .. && pwd )"
pkg_path="github.com/felixge/godrone/cmd"
bin_name="godrone"
drone_ip="${1:-192.168.1.1}"

echo "--> Fetching dependencies ..."
go get "${pkg_path}"
go get "github.com/felixge/makefs"

echo "--> Building http files ..."
go run "${dir}/http_files.go"

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
