#!/bin/sh
set -e

REF = $(cat /github/workflow/event.json | jq '.ref')
FULL_NAME = $(cat /github/workflow/event.json | jq '.repository.full_name')

if [ -z "${KUBECONFIG}" ]; then
  echo "Please set the KUBECONFIG"
  exit 1
fi

if [ -z "${LINK}" ]; then
  LINK= $(cat /github/workflow/event.json | jq '.compare')
fi

if [ -z "${NAME}" ]; then
  NAME= "${FULL_NAME} / ${REF}"
fi

/validate
