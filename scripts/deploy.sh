#!/usr/bin/env bash
set -eu

root_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
bin_name="${1:-godrone}"
drone_ip="${2:-192.168.1.1}"

echo "--> Compiling arm binary ..."
env \
  GOOS=linux \
  GOARCH=arm \
  CGO_ENABLED=0 \
  go install "godrone/cmd/${bin_name}"

echo "--> Uploading via ftp ..."
curl -T "${root_dir}/bin/linux_arm/${bin_name}" "ftp://@${drone_ip}/${bin_name}.next"

echo "--> Starting godrone ..."
"${root_dir}/scripts/start.expect" "${drone_ip}" "${bin_name}"
