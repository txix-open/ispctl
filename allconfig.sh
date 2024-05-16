#!/bin/bash
mkdir -p configs
if [ -n "$*" ]; then
  statusResponse=$(ispctl $* status)
else
  statusResponse=$(ispctl status)
fi

i=0
IFS=$'\n'
for status in $statusResponse; do
  if [ $i -ne 2 ]; then
    i=$((i + 1))
  else
    IFS=$'  '
    read -ra ADDR <<<"$status"
    for module_name in "${ADDR[0]}"; do
      if [ -n "$*" ]; then
        ispctl $* get $module_name . > configs/$module_name.json
      else
        ispctl get $module_name . > configs/$module_name.json
      fi
      break
    done
  fi
done
