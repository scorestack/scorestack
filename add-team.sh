#!/bin/bash

ELASTICSEARCH_HOST=localhost:9200
KIBANA_HOST=localhost:5601
CHECK_FOLDER=examples
USERNAME=root
PASSWORD=changeme

# Wait for elasticsearch to come up
while [[ "$(curl -sku ${USERNAME}:${PASSWORD} "https://${ELASTICSEARCH_HOST}/_cluster/health" | jq -r .status 2>/dev/null)" != "green" ]]
do
  echo "Waiting for Elasticsearch to be ready..."
  sleep 5
done

# Wait for kibana to come up
while [[ "$(curl -sku ${USERNAME}:${PASSWORD} https://${KIBANA_HOST}/api/status | jq -r .status.overall.state 2>/dev/null)" != "green" ]]
do
  echo "Waiting for Kibana to be ready..."
  sleep 5
done

# Add scoreboard dashboard
curl -ku ${USERNAME}:${PASSWORD} https://${KIBANA_HOST}/api/kibana/dashboards/import -H "Content-Type: application/json" -H "kbn-xsrf: true" -d @dashboards/scoreboard.json
curl -kX POST -u ${USERNAME}:${PASSWORD} https://${KIBANA_HOST}/api/spaces/_copy_saved_objects -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"spaces":["scorestack"],"objects":[{"type":"dashboard","id":"'scorestack-scoreboard'"}],"includeReferences":true}'

# Add default index template
curl -k -XPUT -u ${USERNAME}:${PASSWORD} https://${ELASTICSEARCH_HOST}/_template/default -H 'Content-Type: application/json' -d '{"index_patterns":["check*","attrib_*","results*"],"settings":{"number_of_replicas":"0"}}'

# Create results indices
curl -k -XPUT -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/results-admin -H "Content-Type: application/json" -d "@results-admin.json"
curl -k -XPUT -u elastic:${elastic_pass} ${ELASTICSEARCH_HOST}/results-all -H "Content-Type: application/json" -d "@results-all.json"

# Loop through all teams passed as arguments
for TEAM in "${@}"
do
  # Add index for the team results
  curl -k -XPUT -u ${USERNAME}:${PASSWORD} https://${ELASTICSEARCH_HOST}/results-${TEAM} -H "Content-Type: application/json" -d "@results-team.json"

  TEAM_NUM=$(echo $TEAM | sed "s/[a-zA-Z_]//g" | sed "s/^0//g")
  # Add example checks for the team
  for check in $(find ${CHECK_FOLDER} -maxdepth 1 -mindepth 1 -type d -printf "%f\n")
  do
    # Add check definition
    cat ${CHECK_FOLDER}/${check}/check.json | jq --arg TEAM "$TEAM" '.group = $TEAM | .id = "\(.id)-\($TEAM)"' | sed -e "s/\${TEAM_NUM}/${TEAM_NUM}/g" > check.tmp.json
    ID=$(cat check.tmp.json | jq -r '.id')
    curl -k -XPUT -u ${USERNAME}:${PASSWORD} https://${ELASTICSEARCH_HOST}/checkdef/_doc/${ID} -H 'Content-Type: application/json' -d @check.tmp.json
    cat check.tmp.json | jq '{id, name, type, group}' > generic-check.tmp.json
    curl -k -XPUT -u ${USERNAME}:${PASSWORD} https://${ELASTICSEARCH_HOST}/checks/_doc/${ID} -H 'Content-Type: application/json' -d @generic-check.tmp.json

    # Add admin attributes, if they are defined
    if [ -f ${CHECK_FOLDER}/${check}/admin-attribs.json ]
    then
      cat ${CHECK_FOLDER}/${check}/admin-attribs.json | jq --arg TEAM "$TEAM" '.group = $TEAM | .id = "\(.id)-\($TEAM)"' | sed -e "s/\${TEAM_NUM}/${TEAM_NUM}/g" > ${CHECK_FOLDER}/${check}/admin-attribs.tmp.json
      curl -k -XPUT -u ${USERNAME}:${PASSWORD} https://${ELASTICSEARCH_HOST}/attrib_admin_${TEAM}/_doc/${ID} -H "Content-Type: application/json" -d @${CHECK_FOLDER}/${check}/admin-attribs.tmp.json
      rm -f ${CHECK_FOLDER}/${check}/admin-attribs.tmp.json
    fi

    # Add user attributes, if they are defined
    if [ -f ${CHECK_FOLDER}/${check}/user-attribs.json ]
    then
      cat ${CHECK_FOLDER}/${check}/user-attribs.json | jq --arg TEAM "$TEAM" '.group = $TEAM | .id = "\(.id)-\($TEAM)"' | sed -e "s/\${TEAM_NUM}/${TEAM_NUM}/g" > ${CHECK_FOLDER}/${check}/user-attribs.tmp.json
      curl -k -XPUT -u ${USERNAME}:${PASSWORD} https://${ELASTICSEARCH_HOST}/attrib_user_${TEAM}/_doc/${ID} -H "Content-Type: application/json" -d @${CHECK_FOLDER}/${check}/user-attribs.tmp.json
      rm -f ${CHECK_FOLDER}/${check}/user-attribs.tmp.json
    fi
  done

  # Add team role
  curl -kX PUT -u ${USERNAME}:${PASSWORD} https://${KIBANA_HOST}/api/security/role/${TEAM} -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"elasticsearch":{"indices":[{"names":["results-'${TEAM}'*"],"privileges":["read"]},{"names":["attrib_user_'${TEAM}'"],"privileges":["read","index","view_index_metadata"]}]}}'

  # Add team user
  curl -kX PUT -u ${USERNAME}:${PASSWORD} https://${ELASTICSEARCH_HOST}/_security/user/${TEAM} -H 'Content-Type: application/json' -d '{"password":"changeme","roles":["common","'${TEAM}'"]}'

  # Add team overview dashboard
  INDEX="results-${TEAM}"
  CHECKS=$(find examples -maxdepth 1 -mindepth 1 -type d -printf "%f\n" | wc -l)
  cat dashboards/single-team-overview.json | sed -e "s/\${TEAM}/${TEAM}/g" | sed -e "s/\${INDEX}/${INDEX}/g" | sed -e "s/\${CHECKS}/${CHECKS}/g" > tmp-dashboard.json
  curl -ku ${USERNAME}:${PASSWORD} https://${KIBANA_HOST}/api/kibana/dashboards/import -H "Content-Type: application/json" -H "kbn-xsrf: true" -d @tmp-dashboard.json
  curl -kX POST -u ${USERNAME}:${PASSWORD} https://${KIBANA_HOST}/api/spaces/_copy_saved_objects -H 'Content-Type: application/json' -H 'kbn-xsrf: true' -d '{"spaces":["scorestack"],"objects":[{"type":"dashboard","id":"'scorestack-overview-${TEAM}'"}],"includeReferences":true}'
done

# Clean up
rm -f check.tmp.json
rm -f generic-check.tmp.json
rm -f admin-attribs.tmp.json
rm -f tmp-dashboard.json
