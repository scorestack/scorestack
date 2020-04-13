#!/bin/bash

ELASTICSEARCH_HOST=localhost:9200
KIBANA_HOST=localhost:5601
CHECK_FOLDER=examples

# Wait for elasticsearch to come up
while [[ "$(curl -sku root:changeme "https://${ELASTICSEARCH_HOST}/_cluster/health" | jq -r .status 2>/dev/null)" != "green" ]]
do
  echo "Waiting for Elasticsearch to be ready..."
  sleep 5
done

# Wait for kibana to come up
while [[ "$(curl -sku root:changeme https://${KIBANA_HOST}/api/status | jq -r .status.overall.state 2>/dev/null)" != "green" ]]
do
  echo "Waiting for Kibana to be ready..."
  sleep 5
done

# Add scoreboard dashboard
UUID_A=aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
UUID_B=$(uuidgen)
UUID_C=$(uuidgen)
UUID_D=dddddddd-dddd-dddd-dddd-dddddddddddd
UUID_E=eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee
UUID_F=ffffffff-ffff-ffff-ffff-ffffffffffff
cat dashboards/scoreboard.json | sed -e "s/\${UUID_A}/${UUID_A}/g" | sed -e "s/\${UUID_B}/${UUID_B}/g" | sed -e "s/\${UUID_C}/${UUID_C}/g" | sed -e "s/\${UUID_D}/${UUID_D}/g" | sed -e "s/\${UUID_E}/${UUID_E}/g" | sed -e "s/\${UUID_F}/${UUID_F}/g" > tmp-dashboard.json
curl -ku root:changeme https://${KIBANA_HOST}/api/kibana/dashboards/import -H "Content-Type: application/json" -H "kbn-xsrf: true" -d @tmp-dashboard.json
curl -kX POST -u root:changeme https://${KIBANA_HOST}/api/spaces/_copy_saved_objects -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"spaces":["scorestack"],"objects":[{"type":"dashboard","id":"'${UUID_A}'"}],"includeReferences":true}'

# Clean up
rm tmp-dashboard.json

# Add default index template
curl -k -XPUT -u root:changeme https://${ELASTICSEARCH_HOST}/_template/default -H 'Content-Type: application/json' -d '{"index_patterns":["check*","attrib_*"],"settings":{"number_of_replicas":"0"}}'

# Loop through all teams passed as arguments
for TEAM in "${@}"
do
  TEAM_NUM=$(echo $TEAM | sed "s/[a-zA-Z_]//g" | sed "s/^0//g")
  # Add example checks for the team
  for check in $(find ${CHECK_FOLDER} -maxdepth 1 -mindepth 1 -type d -printf "%f\n")
  do
    # Add check definition
    cat ${CHECK_FOLDER}/${check}/check.json | jq --arg TEAM "$TEAM" '.group = $TEAM | .id = "\(.id)-\($TEAM)"' > check.tmp.json
    ID=$(cat check.tmp.json | jq -r '.id')
    curl -k -XPUT -u root:changeme https://${ELASTICSEARCH_HOST}/checkdef/_doc/${ID} -H 'Content-Type: application/json' -d @check.tmp.json
    cat check.tmp.json | jq '{id, name, type, group}' > generic-check.tmp.json
    curl -k -XPUT -u root:changeme https://${ELASTICSEARCH_HOST}/checks/_doc/${ID} -H 'Content-Type: application/json' -d @generic-check.tmp.json

    # Add admin attributes, if they are defined
    if [ -f ${CHECK_FOLDER}/${check}/admin-attribs.json ]
    then
      curl -k -XPUT -u root:changeme https://${ELASTICSEARCH_HOST}/attrib_admin_${TEAM}/_doc/${ID} -H "Content-Type: application/json" -d @${CHECK_FOLDER}/${check}/admin-attribs.json
    fi

    # Add user attributes, if they are defined
    if [ -f ${CHECK_FOLDER}/${check}/user-attribs.json ]
    then
      curl -k -XPUT -u root:changeme https://${ELASTICSEARCH_HOST}/attrib_user_${TEAM}/_doc/${ID} -H "Content-Type: application/json" -d @${CHECK_FOLDER}/${check}/user-attribs.json
    fi
  done

  # Add team role
  curl -kX PUT -u root:changeme https://${KIBANA_HOST}/api/security/role/${TEAM} -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["results-'${TEAM}'*"],"privileges":["read"]},{"names":["attrib_user_*-'${TEAM}'"],"privileges":["read","index","view_index_metadata"]}]}}'

  # Add team user
  curl -kX PUT -u root:changeme https://${ELASTICSEARCH_HOST}/_security/user/${TEAM} -H 'Content-Type: application/json' -d '{"password":"changeme","roles":["common","'${TEAM}'"]}'

  # Add team overview dashboard
  UUID_A=$(uuidgen)
  UUID_B=$(uuidgen)
  UUID_C=$(uuidgen)
  UUID_D=$(uuidgen)
  UUID_E=$(uuidgen)
  UUID_F=$(uuidgen)
  INDEX="results-${TEAM}*"
  CHECKS=$(find examples -maxdepth 1 -mindepth 1 -type d -printf "%f\n" | wc -l)
  cat dashboards/single-team-overview.json | sed -e "s/\${UUID_A}/${UUID_A}/g" | sed -e "s/\${UUID_B}/${UUID_B}/g" | sed -e "s/\${UUID_C}/${UUID_C}/g" | sed -e "s/\${UUID_D}/${UUID_D}/g" | sed -e "s/\${UUID_E}/${UUID_E}/g" | sed -e "s/\${UUID_F}/${UUID_F}/g" | sed -e "s/\${TEAM}/${TEAM}/g" | sed -e "s/\${INDEX}/${INDEX}/g" | sed -e "s/\${CHECKS}/${CHECKS}/g" > tmp-dashboard.json
  curl -ku root:changeme https://${KIBANA_HOST}/api/kibana/dashboards/import -H "Content-Type: application/json" -H "kbn-xsrf: true" -d @tmp-dashboard.json
  curl -kX POST -u root:changeme https://${KIBANA_HOST}/api/spaces/_copy_saved_objects -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"spaces":["scorestack"],"objects":[{"type":"dashboard","id":"'${UUID_A}'"}],"includeReferences":true}'
done

# Clean up
rm -f check.tmp.json
rm -f generic-check.tmp.json
rm -f admin-attribs.tmp.json
rm -f tmp-dashboard.json