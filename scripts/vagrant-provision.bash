#!/usr/bin/env bash

# vagrant-provision.bash provisions an ubuntu linux box with everything needed
# for godrone development. It aims to be idempotent and can continue where it
# left off if interrupted or an error occured.

set -eu
go_version="1.2"
packages="build-essential git-core mercurial"

godrone_dir="$(pwd)"
log_file="${godrone_dir}/vagrant-provision.log"
go_path="${godrone_dir}"
install_dir="${HOME}"
profile_file="${HOME}/.profile"

log() {
  echo "========== $@ ==========" >> "${log_file}"
  echo $@
}

on_error() {
  echo "Failed to provision. See logfile: ${log_file}"
}

profile_add() {
  local line="$1"
  if ! grep "${line}" "${profile_file}" &>> "${log_file}"; then
    log "Adding ${line} to ${profile_file}"
    echo "${line}" >> "${profile_file}"
    eval "${line}"
  fi
}

install_packages() {
  local installed="${install_dir}/packages.txt"

  if [[ "$(cat "${installed}" 2>> "${log_file}")" != "$(echo ${packages})" ]]; then
    log "Updating apt-get"
    sudo apt-get -y update &>> "${log_file}"
    log "Installing apt-get packages: ${packages}"
    sudo apt-get -y install ${packages}  &>> "${log_file}"
    echo "${packages}" > "${installed}"
  fi
}

install_go() {
  local file="go${go_version}.linux-amd64.tar.gz"
  local url="https://go.googlecode.com/files/${file}"
  local dst="/usr/local/go"
  local dst_bin="${dst}/bin"

  pushd "${install_dir}" &>> "${log_file}"
  if [ ! -f "${file}" ]; then
    log "Downloading go ${go_version}"
    curl -s -L -o "${file}" "${url}" &>> "${log_file}"
    sudo rm -rf "${dst}"
  fi

  if [ ! -d "${dst}" ]; then
    log "Extracting go ${go_version}"
    sudo tar -C "$(dirname "${dst}")" -xzf "${file}"

    log "Building go for linux/arm cross-compilation"
    pushd "${dst}/src" &>> "${log_file}"
    sudo env GOOS=linux GOARCH=arm ./make.bash --no-clean &>> "${log_file}"
    popd &>> "${log_file}"
  fi

  profile_add 'export GOPATH="'"${go_path}"'"'
  profile_add 'export PATH="'"${dst_bin}"':${GOPATH}/bin:${PATH}"'

  popd &>> "${log_file}"
}

main() {
  rm -rf "${log_file}"

  local symlink_target="src/github.com/felixge/godrone"
  mkdir -p "$(dirname "${symlink_target}")"
  rm -f "${symlink_target}"
  ln -s ./../../.. "${symlink_target}"

  install_packages
  install_go

  profile_add 'cd "'"${godrone_dir}"'"'
  sudo chown -R "${USER}:${USER}" "${HOME}"
}

trap on_error ERR # @TODO this trap is not working for some reason yet
main
