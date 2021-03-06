#!/bin/sh

# Set up example checks
for check in $(find examples -maxdepth 1 -mindepth 1 -type d -printf "%f\n")
do
  # Add check definition
  curl -k -XPUT -u root:changeme https://localhost:9200/checks/_doc/${check}-example -H 'Content-Type: application/json' -d @examples/${check}/check.json

  # Add admin attributes, if they are defined
  if [ -f examples/${check}/admin-attribs.json ]
  then
    curl -k -XPUT -u root:changeme https://localhost:9200/attrib_admin_${check}-example/_doc/attributes -H "Content-Type: application/json" -d @examples/${check}/admin-attribs.json
  fi

  # Add user attributes, if they are defined
  if [ -f examples/${check}/user-attribs.json ]
  then
    curl -k -XPUT -u root:changeme https://localhost:9200/attrib_user_${check}-example/_doc/attributes -H "Content-Type: application/json" -d @examples/${check}/user-attribs.json
  fi
done


