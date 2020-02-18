#!/bin/sh
until curl -kX GET "https://localhost:9200/_cat/nodes?v&pretty"
do
  sleep 5
done

# This line is really gross but I don't know how else to wrap it. Suggestions
# are welcome!
NEW_PASSWDS=$(docker exec elas01 /bin/bash -c "bin/elasticsearch-setup-passwords auto --batch -Expack.security.http.ssl.key=/usr/share/elasticsearch/config/pki/elas.key.pem -Expack.security.http.ssl.certificate=/usr/share/elasticsearch/config/pki/elas.cert.pem -Expack.security.http.ssl.certificate_authorities=/usr/share/elasticsearch/config/pki/ca.cert.pem --url https://elas01:9200" | grep PASSWORD)

# Extract passwords from utility output
kibana_pass=$(echo "${NEW_PASSWDS}" | grep kibana | awk '{print $NF}')
elastic_pass=$(echo "${NEW_PASSWDS}" | grep elastic | awk '{print $NF}')
beats_pass=$(echo "${NEW_PASSWDS}" | grep beats_system | awk '{print $NF}')
logstash_system_pass=$(echo "${NEW_PASSWDS}" | grep logstash_system | awk '{print $NF}')
logstash_user_pass=$(openssl rand -hex 20)

# Write kibana password to default environment file
cat > .env << EOF
KIBANA_PASSWORD=${kibana_pass}
BEATS_PASSWORD=${beats_pass}
LOGSTASH_USER_PASSWORD=${logstash_user_pass}
LOGSTASH_SYSTEM_PASSWORD=${logstash_system_pass}
EOF

# Install kibana plugin
docker exec kiba01 /bin/bash -c "bin/kibana-plugin install https://tinyurl.com/scorestack-kibana-plugin-zip"

# Create admin user
curl -k -XPOST -u elastic:${elastic_pass} 'https://localhost:9200/_security/user/root' -H "Content-Type: application/json" -d '{"password":"changeme","full_name":"root","email":"root@example.com","roles":["superuser"]}'

# Create logstash user
curl -k -XPOST -u elastic:${elastic_pass} 'https://localhost:9200/_security/role/logstash_writer' -H "Content-Type: application/json" -d '{"cluster":["manage_index_templates","monitor","manage_ilm"],"indices":[{"names":["results-*"],"privileges":["write","create","delete","create_index","manage","manage_ilm"]}]}'
curl -k -XPOST -u elastic:${elastic_pass} 'https://localhost:9200/_security/user/logstash_internal' -H "Content-Type: application/json" -d '{"password":"'"${logstash_user_pass}"'","roles":["logstash_writer"],"full_name":"Internal Logstash User"}'

# Restart kibana and logstash to apply password change
docker-compose up -d --force-recreate kiba01
docker-compose up -d --force-recreate logs01

# Set up example checks
for check in $(find examples -maxdepth 1 -mindepth 1 -type d -printf "%f\n")
do
  # Add check definition
  curl -k -XPUT -u elastic:${elastic_pass} https://localhost:9200/checks/_doc/${check}-example -H 'Content-Type: application/json' -d @examples/${check}/check.json

  # Add admin attributes, if they are defined
  if [ -f examples/${check}/admin-attribs.json ]
  then
    curl -k -XPUT -u elastic:${elastic_pass} https://localhost:9200/attrib_admin_${check}-example/_doc/attributes -H "Content-Type: application/json" -d @examples/${check}/admin-attribs.json
  fi

  # Add user attributes, if they are defined
  if [ -f examples/${check}/user-attribs.json ]
  then
    curl -k -XPUT -u elastic:${elastic_pass} https://localhost:9200/attrib_user_${check}-example/_doc/attributes -H "Content-Type: application/json" -d @examples/${check}/user-attribs.json
  fi
done

# Wait for kibana to be up
while [[ "$(curl -sku root:changeme https://localhost:5601/api/status | jq -r .status.overall.state)" != "green" ]]
do
  echo "Waiting for Kibana to be ready..."
  sleep 5
done

# Set up dashboards
UUID_A=$(uuidgen)
UUID_B=$(uuidgen)
UUID_C=$(uuidgen)
UUID_D=$(uuidgen)
UUID_E=$(uuidgen)
UUID_F=$(uuidgen)
TEAM="Example"
INDEX="results-example*"
cat dashboards/single-team-overview.json | sed -e "s/\${UUID_A}/${UUID_A}/g" | sed -e "s/\${UUID_B}/${UUID_B}/g" | sed -e "s/\${UUID_C}/${UUID_C}/g" | sed -e "s/\${UUID_D}/${UUID_D}/g" | sed -e "s/\${UUID_E}/${UUID_E}/g" | sed -e "s/\${UUID_F}/${UUID_F}/g" | sed -e "s/\${TEAM}/${TEAM}/g" | sed -e "s/\${INDEX}/${INDEX}/g" > tmp-dashboard.json
curl -ku root:changeme https://localhost:5601/api/kibana/dashboards/import -H "Content-Type: application/json" -H "kbn-xsrf: true" -d @tmp-dashboard.json
UUID_A=$(uuidgen)
UUID_B=$(uuidgen)
UUID_C=$(uuidgen)
UUID_D=$(uuidgen)
UUID_E=$(uuidgen)
UUID_F=$(uuidgen)
cat dashboards/scoreboard.json | sed -e "s/\${UUID_A}/${UUID_A}/g" | sed -e "s/\${UUID_B}/${UUID_B}/g" | sed -e "s/\${UUID_C}/${UUID_C}/g" | sed -e "s/\${UUID_D}/${UUID_D}/g" | sed -e "s/\${UUID_E}/${UUID_E}/g" | sed -e "s/\${UUID_F}/${UUID_F}/g" > tmp-dashboard.json
curl -ku root:changeme https://localhost:5601/api/kibana/dashboards/import -H "Content-Type: application/json" -H "kbn-xsrf: true" -d @tmp-dashboard.json
rm tmp-dashboard.json