#!/usr/bin/env bash
set -eu

usage() {
  local exitcode="$1"
  if [[ "$exitcode" = 1 ]]; then
    exec 1>&2
  fi
  echo "Usage: $0 [-h] [--tracegc] [-i <ip>] [-e <envargs>] [<cmd>]"
  exit "$exitcode"
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

  log "Uploading to ${ip}..."
  bash -c "${curlcmd}"
}

clean() {
  rm -rf "$1"
}

main() {
  local ip="192.168.1.1"
  local envargs=''
  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      '-h')
        usage 0
        shift
        ;;
      '-i')
        ip="$2"
        shift
        ;;
      '--tracegc')
        envargs="${envargs} GOGCTRACE=1"
        ;;
      '-e')
        envargs="${envargs} $2"
        shift
        ;;
       -*)
         echo -e "unknown option: $1\n" 1>&2
         usage 1
         ;;
       *)
         break
         ;;
    esac
    shift
  done
  readonly ip
  readonly envargs

  if [[ "$#" -gt 1 ]]; then
    echo -e "unexpected argument: $2\n" 1>&2
    usage 1
  fi

  local readonly cmd="${1:-godrone}"
  local readonly scripts_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
  local readonly dir="$( cd "${scripts_dir}"/.. && pwd )"
  local readonly pkg_path="github.com/felixge/godrone/cmd/${cmd}"
  local readonly startup_script='start.sh'

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
