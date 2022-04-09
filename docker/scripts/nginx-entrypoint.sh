#!/bin/bash
# This entrypoint script is used to dynamically configure where the
# Elasticsearch and Kibana backend services are running. What we gain from this
# added complexity is the ability to deploy Elasticsearch or Kibana clusters
# with the same nginx image.
set -eou pipefail

# Space-separated lists of hosts to proxy to. Ex:
# ELASTIC_BACKEND="es01:9200 es02:9200 es03:9200"
# IMPORTANT: **do not** include a scheme! (https:// or http://)
ELASTIC_BACKEND="${ELASTIC_BACKEND-elasticsearch:9200}"
KIBANA_BACKEND="${KIBANA_BACKEND-kibana:5601}"

# https://stackoverflow.com/a/10586169
IFS=" " read -r -a elastic_upstreams <<< "${ELASTIC_BACKEND}"
IFS=" " read -r -a kibana_upstreams <<< "${KIBANA_BACKEND}"

echo "upstream elasticsearch-backend {" > /etc/nginx/conf.d/elasticsearch-upstream.conf
for upstream in "${elastic_upstreams[@]}"; do
  echo "    server ${upstream};" >> /etc/nginx/conf.d/elasticsearch-upstream.conf
done
echo "}" >> /etc/nginx/conf.d/elasticsearch-upstream.conf

echo "upstream kibana-backend {" > /etc/nginx/conf.d/kibana-upstream.conf
for upstream in "${kibana_upstreams[@]}"; do
  echo "    server ${upstream};" >> /etc/nginx/conf.d/kibana-upstream.conf
done
echo "}" >> /etc/nginx/conf.d/kibana-upstream.conf

exec /docker-entrypoint.sh nginx -g "daemon off;"