Appendix A - FAQ
================

### How can I "reset" my Scorestack cluster without redeploying it?

First, stop all running instances of Dynamicbeat. Once they have all exited, use the Discover app in Kibana to wait for all queued check results to finish processing in Logstash and indexing in Elasticsearch.

Once all check results have been indexed, deleted all the Scorestack indices. You can use the Dev Tools app in Kibana to simplify this process. You need execute the following `DELETE` queries:

- `DELETE /check*`
- `DELETE /attrib*`
- `DELETE /result*`

Finally, rerun the `add-team.sh` script with all the teams you want Scorestack to have. Once the script is finished, you can restart Dynamicbeat.