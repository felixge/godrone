#!/usr/bin/env bash
set -eu

# dist.bash <version>
#
# Creates archive files ready for distribution for all selected platforms.

main() {
  local version="$1"
  local scripts_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
  local root_dir="$( cd "${scripts_dir}"/.. && pwd )"
  local targets=(
    'darwin 386 zip'
    'windows 386 zip'
    'linux 386 tar'
    'linux amd64 tar'
    )

  for target in "${targets[@]}"; do
    build $target
  done
}

build() {
  local os="$1"
  local arch="$2"
  local archive="$3"
  local name="godrone-${version}-${os}-${arch}"
  local out_dir="${root_dir}/dist/${name}"

  echo "Building ${name}"
  "${scripts_dir}/build.bash" \
    -os "${os}" \
    -arch "${arch}" \
    -version ${version} \
    "${out_dir}"
  cd "${out_dir}/.."
  case "${archive}" in
    'zip')
      local archive_name="${out_dir}.zip"
      rm -f "${archive_name}"
      zip -q -r "${archive_name}" "${name}"
      ;;
    'tar')
      local archive_name="${out_dir}.tar.gz"
      rm -f "${archive_name}"
      tar -czf "${archive_name}" "${name}"
      ;;
    *)
      echo "unknown archive format: ${archive}"
      exit 1
      ;;
  esac
  rm -rf "${out_dir}"
}

main $@
