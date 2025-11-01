#!/bin/bash
mkdir -p schemas
if [ -n "$*" ]; then
  statusResponse=$(ispctl $* status)
else
  statusResponse=$(ispctl status)
fi

i=0
IFS=$'\n'
for status in $statusResponse; do
  if [ $i -ne 2 ]; then
    i=$(( i + 1 ))
  else
    IFS=$'  '
    read -ra ADDR <<<"$status"
    for module_name in "${ADDR[0]}"; do
        if [ -n "$*" ]; then
          ispctl $* schema -o html $module_name > schemas/$module_name.html
        else
          ispctl schema -o html $module_name > schemas/$module_name.html
        fi
      break
    done
  fi
done
