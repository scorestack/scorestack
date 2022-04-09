#!/bin/bash
set -eou pipefail

OUTFILE=$(mktemp)
curl --silent --cacert config/certs/ca/ca.crt --head --output "${OUTFILE}" https://localhost:5601

if ! grep --quiet 'HTTP/1.1 302 Found' "${OUTFILE}"; then
  echo "Healthcheck failed. Request headers for healthcheck request:"
  cat "${OUTFILE}"
  exit 1
else
  echo "Healthcheck passed"
fi