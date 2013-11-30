#!/bin/sh
# This script kills any previously running control program (e.g. the original
# parrot firmware, godrone, or godrone diagnostic programs) and then starts
# up the desired program.
set -eu

readonly cmd="$1"
readonly envargs="$2"

# avoid parrot firmware from restarting restarts
touch /tmp/.norespawn
# kill parrot firmware and the cmd we're deploying if it's running
killall -9 program.elf "${cmd}" 2> /dev/null || true

cd /data/video

chmod +x "${cmd}"

# Taken from /bin/program.elf.respawner.sh. Not sure why/if it is needed, seems
# to turn the bottom LED red.
gpio 181 -d ho 1

exec env $envargs "./${cmd}"
