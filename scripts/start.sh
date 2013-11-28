#!/bin/sh
# This script kills any previously running control program (e.g. the original
# parrot firmware, godrone, or godrone diagnostic programs) and then starts
# up the desired program.
set -eu

readonly cmd="${1}"

# avoid restarts
touch /tmp/.norespawn
# @todo inject godrone/navboard
killall -9 program.elf godrone navboard "${cmd}" || true

cd /data/video

rm -f "${cmd}"
mv "${cmd}.next" "${cmd}"
chmod +x "${cmd}"

# Taken from /bin/program.elf.respawner.sh. Not sure why/if it is needed, seems
# to turn the bottom LED red.
gpio 181 -d ho 1

GOGCTRACE=1 exec "./${cmd}"
