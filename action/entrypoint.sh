#!/bin/sh
set -e

if [ -z "${KUBECONFIG}" ]; then
  echo "Please set the KUBECONFIG"
  exit 1
fi

if [ -z "${LINK}" ]; then
  export LINK= $(cat /github/workflow/event.json | jq '.compare')
fi

if [ -z "${NAME}" ]; then
  export NAME= $(cat /github/workflow/event.json | jq '.repository.full_name') + $(cat /github/workflow/event.json | jq '.ref')
fi

./validate
