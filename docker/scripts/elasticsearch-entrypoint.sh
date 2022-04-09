#!/bin/bash
set -eou pipefail

while [ ! -f config/certs/elasticsearch/elasticsearch.crt ]; do
  echo "Certificate not found. Waiting 5s for setup container to finish..."
  sleep 5
done

echo "Starting Elasticsearch"
exec /usr/local/bin/docker-entrypoint.sh eswrapper