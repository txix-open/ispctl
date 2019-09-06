#!/bin/bash
mkdir -p schemas
if [ -n "$*" ]; then
  statusResponse=$(isp-ctl $* status)
else
  statusResponse=$(isp-ctl status)
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
          echo $(isp-ctl $* schema $module_name -o html) >schemas/$module_name.html
        else
          echo $(isp-ctl schema $module_name -o html) >schemas/$module_name.html
        fi

      break
    done
  fi
done
