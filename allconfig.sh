#!/bin/bash
mkdir -p configs
if [ -n "$*" ]; then
  statusResponse=$(isp-ctl $* status)
else
  statusResponse=$(isp-ctl status)
fi

firstRecordSkip="MODULE"
secondRecordSkip="+------------------------+---------------+-------------------------------+"
IFS=$'\n'
for status in $statusResponse; do
  IFS=$'  '
  read -ra ADDR <<<"$status"
  for module_name in "${ADDR[0]}"; do
    if [ "$module_name" != "$firstRecordSkip" ] && [ "$module_name" != "$secondRecordSkip" ]; then
      if [ -n "$*" ]; then
        echo $(isp-ctl $* get $module_name .) >configs/$module_name.json
      else
        echo $(isp-ctl get $module_name .) >configs/$module_name.json
      fi
    fi
    break
  done
done