Appendix A - FAQ
================

### How can I "reset" my Scorestack cluster without redeploying it?

First, stop all running instances of Dynamicbeat. Once they have all exited, use the Discover app in Kibana to wait for all queued check results to finish processing in Logstash and indexing in Elasticsearch.

Once all check results have been indexed, deleted all the Scorestack indices. You can use the Dev Tools app in Kibana to simplify this process. You need execute the following `DELETE` queries:

- `DELETE /check*`
- `DELETE /attrib*`
- `DELETE /result*`

Finally, rerun the `add-team.sh` script with all the teams you want Scorestack to have. Once the script is finished, you can restart Dynamicbeat.

### What are the default admin credentials?

The following credentials are the default `superuser` credentials for all Scorestack deployments:

```
Username: root
Password: changeme
```

These credentials will give you the same access as the [`elastic` built-in user](https://www.elastic.co/guide/en/elasticsearch/reference/current/built-in-users.html), so be careful, and make sure to change the password after your first log-in!

### What are the default Dynamicbeat credentials?

Dynamicbeat needs a username and password to connect to Elasticsearch and read the current check definitions and attributes. A user is created for you during all Scorestack deployments with the following default credentials:

```
Username: dynamicbeat
Password: changeme
```