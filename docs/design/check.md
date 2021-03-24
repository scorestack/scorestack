Life of a Check
===============

This document walks through the life of a check, starting at loading definitions from Elasticsearch and ending at indexing check results back into Elasticsearch. This has mostly the same information as the [Elastic Stack Architecture](architecture.md) document, but presented as a timeline rather than an explanation of components.

Running Checks
--------------

A **check** is a single attempt to verify the functionality of a specific network service. Here are a few examples of basic checks that could be run:

- Send a series of ICMP Echo Requests to an IP address and expect a certain number or percentage of ICMP Echo Replies from that address.
- Send an HTTP GET request for a specific webpage to a webserver and check that the returned content matches the expected content.
- Log in to a system via SSH with a given set of credentials, run a specific command, and ensure the command printed the expected content.

Each check is defined by a **check definition**, which is a JSON document stored in Elasticsearch that provides the necessary information for Dynamicbeat to run the check. Additionally, a check may have **check attributes**, which allow for on-the-fly modification of variables that are templated into the check definition. For more information on defining checks, see [the check definition documentation](./checks.md).

When Dynamicbeat first starts, it pulls the check definitions stored in Elasticsearch and stores them in memory. Then every 30 seconds, Dynamicbeat will start a single check for each one of the check definitions that it currently has stored in memory. Additionally, every minute Dynamicbeat will refresh the check definition information that it has stored in memory by querying Elasticsearch for updates.

Reporting Check Results
-----------------------

As checks finish executing, their results (pass, fail, or timeout, with some additional information) are buffered to be sent to Elasticserach later on. Typically, Dynamicbeat will immediately send check results to Elasticsearch. However, if there issues establishing a stable connection to Elasticsearch, the results will stay in Dynamicbeat's buffer until it can reestablish the connection.

Before Dynamicbeat indexes a check result in Elasticsearch, it will perform some basic processing.

First, the `passed` boolean field is converted to an integer in the `passed_int` field. If the check passed, `passed_int` will be set to `1`. Otherwise, `passed_int` will be `0`. This conversion allows for easy score calculation within Kibana dashboards.

Next, the `@timestamp` field is converted to an integer representing the Unix epoch representation of the timestamp, which is stored in the `epoch` field. This conversion makes it simple to display only the latest check results within Kibana dashboards.

Finally, three versions of the result event are created: generic, admin, and group. These events are then stored in an Elasticsearch index that matches the glob `results-*-TIMESTAMP`, where `TIMESTAMP` is a timestamp representing the current date in the format `YYYY.MM.DD`.

Generic Results
---------------

Generic results have the `message` and `details` fields removed, and are viewable by all Scorestack users. This allows teams to see how other teams are doing, but does not give them information on _why_ other teams' checks may be failing. Since field-based access control is a premium feature of the Elastic Stack, this workaround is required for competition-wide dashboards to work without revealing details of check results to other teams.

Generic results are stored in the `results-all-*` indices.

Group Results
-------------

The group results do not have the `message` and `details` fields removed, and are only viewable by members of the group the check belongs to. For example, the group results for a check in the Team 3 group can only be viewed by users who are in the Team 3 group. Group results allow users to get more detailed information about why their checks may be failing, which provides them a starting point for troubleshooting their services.

Group results are stored in the `results-GROUP-*` indices, where `GROUP` is the group ID.

Admin Results
-------------

The admin results, like the group results, do not have the `message` and `details` fields removed. However, admin results are only viewable by members of the `spectator` group. Admin results are mainly useful for troubleshooting issues with service deployment prior to a competition, or for detecting issues with check definitions and/or Dynamicbeat. This is becase all the admin results are stored within a single set of indices, unlike the group events, which are stored across a variety of indices. Having all admin results in a single set of indices allows Scorestack administrators to search across all check results using a single index glob.

The admin results are stored in the `results-admin-*` indices.