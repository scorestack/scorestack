#!/bin/bash

# Install dependencies
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install -y -q -e 0 epel-release
yum install -y -q -e 0 unzip openssl jq docker-ce-cli
curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

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
  -Expack.security.http.ssl.certificate_authorities=/usr/share/elasticsearch/config/certificates/ca/ca.crt \
  --url ${ELASTICSEARCH_HOST}" | grep PASSWORD > /tmp/cluster-passwords.txt

# Extract passwords from output
kibana_pass=$(cat /tmp/cluster-passwords.txt | grep 'kibana =' | awk '{print $NF}')
elastic_pass=$(cat /tmp/cluster-passwords.txt | grep elastic | awk '{print $NF}')
beats_pass=$(cat /tmp/cluster-passwords.txt | grep beats_system | awk '{print $NF}')
logstash_pass=$(openssl rand -hex 20)

# Set passwords in kibana and logstash keystores
docker exec ${KIBANA_CONTAINER} bin/kibana-keystore create
docker exec ${KIBANA_CONTAINER} /bin/bash -c "bin/kibana-keystore add elasticsearch.password --stdin <<< '${kibana_pass}'"
docker exec ${LOGSTASH_CONTAINER} /bin/bash -c "yes y | bin/logstash-keystore create"
while [ $? -ne 0 ]
do
  docker exec ${LOGSTASH_CONTAINER} /bin/bash -c "yes y | bin/logstash-keystore create"
done
docker exec ${LOGSTASH_CONTAINER} /bin/bash -c "bin/logstash-keystore add elasticsearch_password --stdin <<< '${logstash_pass}'"
while [ $? -ne 0 ]
do
  docker exec ${LOGSTASH_CONTAINER} /bin/bash -c "bin/logstash-keystore add elasticsearch_password --stdin <<< '${logstash_pass}'"
done
docker exec ${LOGSTASH_CONTAINER} /bin/bash -c "bin/logstash-keystore list"

# Write passwords to docker-compose default environment file
cat > config/.env << EOF
BEATS_PASSWORD=${beats_pass}
EOF

# Delete the passwords file
shred -uvz /tmp/cluster-passwords.txt

# Install kibana plugin
docker exec ${KIBANA_CONTAINER} /bin/bash -c "bin/kibana-plugin install https://github.com/scorestack/scorestack/releases/download/v0.6.0/kibana-plugin-v0.6.0.zip"

# Create admin user
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/user/root -H "Content-Type: application/json" -d '{"password":"changeme","full_name":"Extra Superuser","email":"root@example.com","roles":["superuser"]}'

# Add dynamicbeat role and user
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/role/dynamicbeat_reader -H "Content-Type: application/json" -d '{"indices":[{"names":["checkdef*","attrib_*"],"privileges":["read"]}]}'
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/user/dynamicbeat -H "Content-Type: application/json" -d '{"password":"changeme","full_name":"Dynamicbeat Definition-Reading User","email":"dynamicbeat@example.com","roles":["dynamicbeat_reader"]}'

# Create logstash user
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/role/logstash_writer -H "Content-Type: application/json" -d '{"cluster":["manage_index_templates","monitor","manage_ilm"],"indices":[{"names":["results-*"],"privileges":["write","create","delete","create_index","manage","manage_ilm"]}]}'
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/user/logstash -H "Content-Type: application/json" -d '{"password":"'"${logstash_pass}"'","roles":["logstash_writer"],"full_name":"Internal Logstash User","email":"logstash@example.com"}'

# Restart kibana and logstash to reload credentials from keystore
cd config
docker-compose -p docker restart kibana logstash
cd ..

# Wait for kibana to be up
while [[ "$(curl -sku root:changeme ${KIBANA_HOST}/api/status | jq -r .status.overall.state 2>/dev/null)" != "green" ]]
do
  echo "Waiting for Kibana to be ready..."
  sleep 5
done

# Add Scorestack space
curl -kX POST -u root:changeme ${KIBANA_HOST}/api/spaces/space -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"id":"scorestack","name":"Scorestack","disabledFeatures":["visualize","dev_tools","indexPatterns","savedObjectsManagement","graph","monitoring","ml","apm","maps","canvas","infrastructure","logs","siem","uptime"]}'

# Set dark theme on both spaces
curl -kX POST -u root:changeme ${KIBANA_HOST}/api/kibana/settings/theme:darkMode -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"value":"true"}'
curl -kX POST -u root:changeme ${KIBANA_HOST}/s/scorestack/api/kibana/settings/theme:darkMode -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"value":"true"}'

# Add base role for common permissions
curl -kX PUT -u root:changeme ${KIBANA_HOST}/api/security/role/common -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["results-all*","checks"],"privileges":["read"]}]},"kibana":[{"base":["read"],"spaces":["scorestack"]}]}'

# Add spectator role
curl -kX PUT -u root:changeme ${KIBANA_HOST}/api/security/role/spectator -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["results*"],"privileges":["read"]}]}}'

# Add admin roles
curl -kX PUT -u root:changeme ${KIBANA_HOST}/api/security/role/attribute-admin -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["attrib_*"],"privileges":["all"]}]}}'
curl -kX PUT -u root:changeme ${KIBANA_HOST}/api/security/role/check-admin -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["check*"],"privileges":["all"]}]}}'
