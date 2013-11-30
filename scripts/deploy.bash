#!/usr/bin/env bash
set -eu

readonly cmd="${1:-godrone}"
readonly scripts_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
readonly dir="$( cd "${scripts_dir}"/.. && pwd )"
readonly pkg_path="github.com/felixge/godrone/cmd/${cmd}"
readonly drone_ip="192.168.1.1"
readonly startup_script="start.sh"
readonly envargs="GOGCTRACE=1"

echo "--> Fetching dependencies ..."
go get "${pkg_path}"

echo "--> Compiling arm binary ..."
env \
  GOOS=linux \
  GOARCH=arm \
  CGO_ENABLED=0 \
  go build -o "${cmd}" "${pkg_path}"

echo "--> Uploading via ftp ..."
curl \
  -T "${scripts_dir}/${startup_script}" "ftp://@${drone_ip}/${startup_script}" \
  -T "${cmd}" "ftp://@${drone_ip}/${cmd}.next"

rm -rf "${cmd}"

echo "--> Starting godrone ..."
"${scripts_dir}/start.expect" \
  "${drone_ip}" \
  "${startup_script}" \
  "${cmd}" \
  "${envargs}"
