#!/bin/sh
curl --cacert /usr/share/elasticsearch/config/pki/ca.cert.pem -s https://localhost:9200 > /dev/null
if [[ $$? == 52]]
then
  echo 0
else
  echo 1
fi