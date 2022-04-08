Appendix A - FAQ
================

### How can I "reset" my Scorestack cluster without redeploying it?

First, stop all running instances of Dynamicbeat. Once they have all exited, use the Discover app in Kibana to wait for all queued check results to finish indexing in Elasticsearch.

Once all check results have been indexed, deleted all the Scorestack indices. You can use the Dev Tools app in Kibana to simplify this process. You need execute the following `DELETE` queries:

- `DELETE /check*`
- `DELETE /attrib*`
- `DELETE /result*`

Alternatively, you can use this bash one liner to clear the indices: 

`for i in "/check*" "/attrib*" "/result*"; do curl -kXDELETE -u elastic:changeme https://localhost:9200$i && echo $i; done`

Finally, setup Kibana and Elastic again using the `dynamicbeat setup` command and you're good to go!

### What are the default admin credentials?

The following credentials are the default `superuser` credentials for all Scorestack deployments:

```
Username: elastic
Password: changeme
```

Be careful as this user has `superuser` privileges, and make sure to change the password after your first log-in!

### What are the default Dynamicbeat credentials?

Dynamicbeat needs a username and password to connect to Elasticsearch and read the current check definitions and attributes. A user is created for you during all Scorestack deployments with the following default credentials:

```
Username: dynamicbeat
Password: changeme
```
