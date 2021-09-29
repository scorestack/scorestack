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

# Set passwords in kibana keystore
docker exec ${KIBANA_CONTAINER} bin/kibana-keystore create
docker exec ${KIBANA_CONTAINER} /bin/bash -c "bin/kibana-keystore add elasticsearch.password --stdin <<< '${kibana_pass}'"

# Install kibana plugin
docker exec ${KIBANA_CONTAINER} /bin/bash -c "bin/kibana-plugin install https://github.com/scorestack/scorestack/releases/download/v0.8.2/kibana-plugin-v0.8.2.zip"

# Restart kibana to reload credentials from keystore
cd config
docker-compose -p docker restart kibana
cd ..

# Write passwords to docker-compose default environment file
cat > config/.env << EOF
BEATS_PASSWORD=${beats_pass}
EOF

# Delete the passwords file
shred -uvz /tmp/cluster-passwords.txt

# Set Elastic admin password
curl -k -XPOST -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/_security/user/elastic/_password -H "Content-Type: application/json" -d '{"password":"changeme"}'