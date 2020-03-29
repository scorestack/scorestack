#!/bin/bash

# Install dependencies
yum install -y -q -e 0 unzip openssl

# Generate certificate bundle if it isn't already generated
if [[ ! -f /certificates/bundle.zip ]]
then
  bin/elasticsearch-certutil cert --silent --pem --in config/certificates/instances.yml -out /certificates/bundle.zip
  unzip /certificates/bundle.zip -d /certificates
fi

# Convert logstash key into PKCS#8 format
openssl pkcs8 -topk8 -nocrypt -in /certificates/logstash/logstash.key -out /certificates/logstash/logstash.key.pkcs8

# Set proper permissions on certificates directory
chown -R 1000:0 /certificates
