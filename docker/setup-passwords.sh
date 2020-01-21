#!/bin/sh
until curl -kX GET "https://localhost:9200/_cat/notes?v&pretty"
do
  sleep 5
done

# This line is really gross but I don't know how else to wrap it. Suggestions
# are welcome!
NEW_PASSWDS=$(docker exec elas01 /bin/bash -c "bin/elasticsearch-setup-passwords auto --batch -Expack.security.http.ssl.key=/usr/share/elasticsearch/config/pki/elas.key.pem -Expack.security.http.ssl.certificate=/usr/share/elasticsearch/config/pki/elas.cert.pem -Expack.security.http.ssl.certificate_authorities=/usr/share/elasticsearch/config/pki/ca.cert.pem --url https://elas01:9200" | grep PASSWORD)

# Extract passwords from utility output
apm_system_pass=$(echo ${NEW_PASSWDS} | grep apm_system | awk '{print $NF}')
kibana_pass=$(echo ${NEW_PASSWDS} | grep kibana | awk '{print $NF}')
logstash_system_pass=$(echo ${NEW_PASSWDS} | grep logstash_system | awk '{print $NF}')
beats_system_pass=$(echo ${NEW_PASSWDS} | grep beats_system | awk '{print $NF}')
remote_monitoring_user_pass=$(echo ${NEW_PASSWDS} | grep remote_monitoring_user | awk '{print $NF}')
elastic_pass=$(echo ${NEW_PASSWDS} | grep elastic | awk '{print $NF}')