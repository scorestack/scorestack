#!/bin/bash
set -eou pipefail

while [ ! -f config/certs/kibana/kibana.crt ]; do
  echo "Certificate not found. Waiting 5s for setup container to finish..."
  sleep 5
done

while [ ! -f /opt/plugin/build/scorestack-8.1.2.zip ]; do
  echo "Plugin zipfile not found. Waiting 5s for plugin-builder container to finish..."
  sleep 5
done

echo "Installing Kibana plugin"
bin/kibana-plugin install file:///opt/plugin/build/scorestack-8.1.2.zip

echo "Starting Kibana"
exec /usr/local/bin/kibana-docker