#!/usr/bin/env bash
set -eu

scripts_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pkg_path="github.com/felixge/godrone/cmd"
bin_name="godrone"
script_name="${bin_name}.sh"
drone_ip="${1:-192.168.1.1}"

echo "--> Fetching dependencies ..."
go get "${pkg_path}"
go get "github.com/felixge/makefs"

echo "--> Compiling arm binary ..."
env \
  GOOS=linux \
  GOARCH=arm \
  CGO_ENABLED=0 \
  go build -o "${scripts_dir}/${bin_name}" "${pkg_path}"

echo "--> Uploading via ftp ..."
curl \
  -T "${scripts_dir}/${script_name}" "ftp://@${drone_ip}/${script_name}" \
  -T "${scripts_dir}/${bin_name}" "ftp://@${drone_ip}/${bin_name}.next"

echo "--> Starting godrone ..."
"${scripts_dir}/deploy.expect" "${drone_ip}" "${bin_name}" "${script_name}"
