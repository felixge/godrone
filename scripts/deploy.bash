#!/usr/bin/env bash
set -eu

usage() {
  echo "Usage: $0"
}

log() {
  echo "--> $@"
}

fetch_deps() {
  log 'Fetching dependencies ...'
  go get "$1"
}

build() {
  log "Building $2 arm binary ..."
  env \
    GOOS=linux \
    GOARCH=arm \
    CGO_ENABLED=0 \
    go build -o "$2" "$1"
}

upload() {
  local ip="$1"
  shift

  local curlcmd="curl"
  for file in $@; do
    curlcmd=" ${curlcmd} -T '${file}' 'ftp://@${ip}/$(basename "${file}")'"
  done

  log 'Uploading via ftp ...'
  bash -c "${curlcmd}"
}

clean() {
  rm -rf "$1"
}

main() {
  local readonly cmd="${1:-godrone}"
  local readonly scripts_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
  local readonly dir="$( cd "${scripts_dir}"/.. && pwd )"
  local readonly pkg_path="github.com/felixge/godrone/cmd/${cmd}"
  local readonly ip='192.168.1.1'
  local readonly startup_script='start.sh'
  local readonly envargs='GOGCTRACE=1'

  fetch_deps "${pkg_path}"
  build "${pkg_path}" "${cmd}"
  upload "${ip}" "${scripts_dir}/${startup_script}" "${cmd}"
  clean "${cmd}"

  log "Starting ${cmd} ..."
  "${scripts_dir}/start.expect" \
    "${ip}" \
    "${startup_script}" \
    "${cmd}" \
    "${envargs}"
}

main $@
