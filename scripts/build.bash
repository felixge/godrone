#!/usr/bin/env bash
set -eu

main() {
  local scripts_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
  local root_dir="$( cd "${scripts_dir}"/.. && pwd )"
  local base_pkg='github.com/felixge/godrone'
  local arch="$(go env GOARCH)"
  local os="$(go env GOOS)"
  local ldflags=''
  local deploy_name='deploy'

  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      '-arch')
        arch="$2"
        shift
        ;;
      '-os')
        os="$2"
        shift
        ;;
      '-version')
        ldflags="-X main.Version $2"
        shift
        ;;
       *)
         break
         ;;
    esac
    shift
  done
  local out_dir="$1"

  if [[ "${os}" = "windows" ]]; then
    deploy_name="${deploy_name}.exe"
  fi

  rm -rf "${out_dir}"
  mkdir -p "${out_dir}"

  makefs "${root_dir}/http/fs"
  env \
    GOOS="${os}" \
    GOARCH="${arch}" \
    go build -o "${out_dir}/${deploy_name}" "${base_pkg}/cmd/deploy"
  env \
    GOOS=linux \
    GOARCH=arm \
    GOARM=7 \
    CGO_ENABLED=0 \
    go build \
      -o "${out_dir}/godrone" \
      -ldflags "${ldflags}" \
      "${base_pkg}/cmd/godrone"
  cp "${scripts_dir}/start.sh" "${out_dir}"
  cp "${root_dir}/LICENSE.txt" "${out_dir}"
  cp "${root_dir}/cmd/godrone/godrone.conf" "${out_dir}"
}

main $@
