#!/bin/bash

# Install dependencies
yum install -y -q -e 0 unzip openssl jq

# Generate certificate bundle if it isn't already generated
if [[ ! -f /certificates/bundle.zip ]]
then
  bin/elasticsearch-certutil cert --silent --pem --in config/instances.yml -out /certificates/bundle.zip
  unzip /certificates/bundle.zip -d /certificates
fi

# Convert logstash key into PKCS#8 format
openssl pkcs8 -topk8 -nocrypt -in /certificates/logstash/logstash.key -out /certificates/logstash/logstash.key.pkcs8

# Set proper permissions on certificates directory
chown -R 1000:0 /certificates

# Wait for elasticsearch to come up
until curl -kX GET "${ELASTICSEARCH_HOST}/_cat/nodes?v&pretty"
do
  sleep 5
done

# Generate passwords
docker exec ${ELASTICSEARCH_CONTAINER} /bin/bash -c \
  "bin/elasticsearch-setup-passwords auto --batch \
  -Expack.security.http.ssl.key=/usr/share/elasticsearch/config/certificates/elasticsearch/elasticsearch.key \
  -Expack.security.http.ssl.certificate=/usr/share/elasticsearch/config/certificates/elasticsearch/elasticsearch.crt \
  -Expack.security.http.ssl.certificate_authorities=/usr/share/elasticsearch/config/certificates/ca.crt \
  --url ${ELASTICSEARCH_HOST}" | grep PASSWORD > /tmp/cluster-passwords.txt

# Extract passwords from output
kibana_pass=$(cat /tmp/cluster-passwords.txt | grep kibana | awk '{print $NF}')
elastic_pass=$(cat /tmp/cluster-passwords.txt | grep elastic | awk '{print $NF}')
beats_pass=$(cat /tmp/cluster-passwords.txt | grep beats_system | awk '{print $NF}')
logstash_system_pass=$(cat /tmp/cluster-passwords.txt | grep logstash_system | awk '{print $NF}')
logstash_user_pass=$(openssl rand -hex 20)

# Write passwords to docker-compose default environment file
cat > config/.env << EOF
KIBANA_PASSWORD=${kibana_pass}
BEATS_PASSWORD=${beats_pass}
LOGSTASH_USER_PASSWORD=${logstash_user_pass}
LOGSTASH_SYSTEM_PASSWORD=${logstash_system_pass}
EOF

# Delete the passwords file
shred -uvz /tmp/cluster-passwords.txt

# Install kibana plugin
docker exec ${KIBANA_CONTAINER} /bin/bash -c "bin/kibana-plugin install https://tinyurl.com/scorestack-kibana-plugin"

# Create admin user
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/user/root -H "Content-Type: application/json" -d '{"password":"changeme","full_name":"root","email":"root@example.com","roles":["superuser"]}'

# Add dynamicbeat role and user 
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/role/dynamicbeat-role -H "Content-Type: application/json" -d '{"indices":[{"names":["checkdef*","attrib_*"],"privileges":["read"]}]}'
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/user/dynamicbeat -H "Content-Type: application/json" -d '{"password":"changeme","full_name":"dynamicbeat","email":"dynamicbeat@example.com","roles":["dynamicbeat-role"]}'

# Create logstash user
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/role/logstash_writer -H "Content-Type: application/json" -d '{"cluster":["manage_index_templates","monitor","manage_ilm"],"indices":[{"names":["results-*"],"privileges":["write","create","delete","create_index","manage","manage_ilm"]}]}'
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/user/logstash_internal -H "Content-Type: application/json" -d '{"password":"'"${logstash_user_pass}"'","roles":["logstash_writer"],"full_name":"Internal Logstash User"}'

# Wait for kibana to be up
while [[ "$(curl -sku root:changeme ${KIBANA_HOST}/api/status | jq -r .status.overall.state 2>/dev/null)" != "green" ]]
do
  echo "Waiting for Kibana to be ready..."
  sleep 5
done

# Add ScoreStack space
curl -kX POST -u root:changeme ${KIBANA_HOST}/api/spaces/space -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"id":"scorestack","name":"ScoreStack","disabledFeatures":["visualize","dev_tools","advancedSettings","indexPatterns","savedObjectsManagement","graph","monitoring","ml","apm","maps","canvas","infrastructure","logs","siem","uptime"]}'

# Add base role for common permissions
curl -kX PUT -u root:changeme ${KIBANA_HOST}/api/security/role/common -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["results-all*","checks"],"privileges":["read"]}]},"kibana":[{"base":["read"],"spaces":["scorestack"]}]}'

# Add spectator role
curl -kX PUT -u root:changeme ${KIBANA_HOST}/api/security/role/spectator -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["results*"],"privileges":["read"]}]}}'

# Add admin roles
curl -kX PUT -u root:changeme ${KIBANA_HOST}/api/security/role/attribute-admin -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["attrib_*"],"privileges":["all"]}]}}'
curl -kX PUT -u root:changeme ${KIBANA_HOST}/api/security/role/check-admin -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["check*"],"privileges":["all"]}]}}'