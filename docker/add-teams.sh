#!/bin/bash

#for i in {1..5}
#do
  #TEAM="team${i}"
  TEAM="example"

  # Add example checks for the team
  for check in $(find examples -maxdepth 1 -mindepth 1 -type d -printf "%f\n")
  do
    # Add check definition
    cat examples/${check}/check.json | jq --arg TEAM "$TEAM" '.group = $TEAM | .id = "\(.id)-\($TEAM)"' > check.tmp.json
    curl -k -XPUT -u root:changeme https://localhost:9200/checks/_doc/${check}-${TEAM} -H 'Content-Type: application/json' -d @check.tmp.json
    rm check.tmp.json

    # Add admin attributes, if they are defined
    if [ -f examples/${check}/admin-attribs.json ]
    then
      curl -k -XPUT -u root:changeme https://localhost:9200/attrib_admin_${check}-${TEAM}/_doc/attributes -H "Content-Type: application/json" -d @examples/${check}/admin-attribs.json
    fi

    # Add user attributes, if they are defined
    if [ -f examples/${check}/user-attribs.json ]
    then
      curl -k -XPUT -u root:changeme https://localhost:9200/attrib_user_${check}-${TEAM}/_doc/attributes -H "Content-Type: application/json" -d @examples/${check}/user-attribs.json
    fi
  done

  # Add team role

  # Add team user

  # Add team dashboards
  UUID_A=$(uuidgen)
  UUID_B=$(uuidgen)
  UUID_C=$(uuidgen)
  UUID_D=$(uuidgen)
  UUID_E=$(uuidgen)
  UUID_F=$(uuidgen)
  INDEX="results-${TEAM}*"
  CHECKS=$(find examples -maxdepth 1 -mindepth 1 -type d -printf "%f\n" | wc -l)
  cat dashboards/single-team-overview.json | sed -e "s/\${UUID_A}/${UUID_A}/g" | sed -e "s/\${UUID_B}/${UUID_B}/g" | sed -e "s/\${UUID_C}/${UUID_C}/g" | sed -e "s/\${UUID_D}/${UUID_D}/g" | sed -e "s/\${UUID_E}/${UUID_E}/g" | sed -e "s/\${UUID_F}/${UUID_F}/g" | sed -e "s/\${TEAM}/${TEAM}/g" | sed -e "s/\${INDEX}/${INDEX}/g" > tmp-dashboard.json
  curl -ku root:changeme https://localhost:5601/api/kibana/dashboards/import -H "Content-Type: application/json" -H "kbn-xsrf: true" -d @tmp-dashboard.json
  rm tmp-dashboard.json
#done
