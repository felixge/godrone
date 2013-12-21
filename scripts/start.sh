#!/bin/sh
set -eu

src_dir="$(pwd)"
target_dir="${src_dir/.next/}"

rm -rf "${target_dir}"
mv "${src_dir}" "${target_dir}"
cd "${target_dir}"

# avoid parrot firmware getting restarted after killing it
touch /tmp/.norespawn
# kill parrot firmware and godrone
killall -9 program.elf godrone 2> /dev/null || true

# Taken from /bin/program.elf.respawner.sh. Not sure why/if it is needed, seems
# to turn the bottom LED red.
gpio 181 -d ho 1

chmod +x godrone
./godrone
