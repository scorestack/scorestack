#!/bin/bash
set -eou pipefail

if [ ! -f config/certs/ca.zip ]; then
  echo "Creating CA"
  mkdir -p config/certs
  bin/elasticsearch-certutil ca --silent --pem --out config/certs/ca.zip
  unzip config/certs/ca.zip -d config/certs
fi

if [ ! -f config/certs/certs.zip ]; then
  echo "Creating certs"
  bin/elasticsearch-certutil cert --silent --pem \
    --out config/certs/certs.zip \
    --in config/certs/instances.yml \
    --ca-cert config/certs/ca/ca.crt \
    --ca-key config/certs/ca/ca.key
  unzip config/certs/certs.zip -d config/certs
fi

echo "Waiting for Elasticsearch availability"
until curl -s --cacert config/certs/ca/ca.crt https://elasticsearch:9200 | grep -q "missing authentication credentials"; do
  sleep 5
done

echo "Setting kibana_system password"
until curl -s -X POST --cacert config/certs/ca/ca.crt -u "elastic:${ELASTIC_PASSWORD}" -H "Content-Type: application/json" https://elasticsearch:9200/_security/user/kibana_system/_password -d "{\"password\":\"${KIBANA_PASSWORD}\"}" | grep -q "^{}"; do
  sleep 5
done

echo "All done!"