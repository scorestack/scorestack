#!/bin/bash

# Loop through all teams passed as arguments
for TEAM in "${@}"
do
  TEAM_NUM=$(echo $TEAM | sed "s/[a-zA-Z_]//g" | sed "s/^0//g")
  # Add example checks for the team
  for check in $(find examples -maxdepth 1 -mindepth 1 -type d -printf "%f\n")
  do
    # Add check definition
    cat examples/${check}/check.json | jq --arg TEAM "$TEAM" '.group = $TEAM | .id = "\(.id)-\($TEAM)"' > check.tmp.json
    ID=$(cat check.tmp.json | jq -r '.id')
    curl -k -XPUT -u root:changeme https://localhost:9200/checkdef/_doc/${ID} -H 'Content-Type: application/json' -d @check.tmp.json
    cat check.tmp.json | jq '{id, name, type, group}' > generic-check.tmp.json
    curl -k -XPUT -u root:changeme https://localhost:9200/checks/_doc/${ID} -H 'Content-Type: application/json' -d @generic-check.tmp.json

    # Add admin attributes, if they are defined
    if [ -f examples/${check}/admin-attribs.json ]
    then
      curl -k -XPUT -u root:changeme https://localhost:9200/attrib_admin_${ID}/_doc/attributes -H "Content-Type: application/json" -d @examples/${check}/admin-attribs.json
    fi

    # Add user attributes, if they are defined
    if [ -f examples/${check}/user-attribs.json ]
    then
      curl -k -XPUT -u root:changeme https://localhost:9200/attrib_user_${ID}/_doc/attributes -H "Content-Type: application/json" -d @examples/${check}/user-attribs.json
    fi
  done

  # Add team role
  curl -kX PUT -u root:changeme https://localhost:5601/api/security/role/${TEAM} -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["results-'${TEAM}'*"],"privileges":["read"]},{"names":["attrib_user_*-'${TEAM}'"],"privileges":["read","index","view_index_metadata"]}]}}'

  # Add team user
  curl -kX PUT -u root:changeme https://localhost:9200/_security/user/${TEAM} -H 'Content-Type: application/json' -d '{"password":"changeme","roles":["common","'${TEAM}'"]}'

  # Wait for kibana to be up
  while [[ "$(curl -sku root:changeme https://localhost:5601/api/status | jq -r .status.overall.state 2>/dev/null)" != "green" ]]
  do
    echo "Waiting for Kibana to be ready..."
    sleep 5
  done

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
  curl -ku root:changeme https://localhost:5601/api/kibana/dashboards/import -H "Content-Type: application/json" -H "kbn-xsrf: true" -d @tmp-dashboard.json
  curl -kX POST -u root:changeme https://localhost:5601/api/spaces/_copy_saved_objects -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"spaces":["scorestack"],"objects":[{"type":"dashboard","id":"'${UUID_A}'"}],"includeReferences":true}'
done

# Clean up
rm -f check.tmp.json
rm -f generic-check.tmp.json
rm -f admin-attribs.tmp.json
rm -f tmp-dashboard.json