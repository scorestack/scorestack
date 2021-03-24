Elastic Stack Architectural Overview
====================================

Scorestack is based around a customized deployment of the Elastic Stack that includes Elasticsearch, Kibana, a Kibana plugin, and a service-checking program named Dynamicbeat. All of these components are configured with X-Pack security enabled, and TLS mutual authentication is used for all inter-cluster communications.

This document provides an overview of each of the Elastic Stack components used in Scorestack, what they are used for, and what they do. It has mostly the same information as the [Life of a Check](check.md) document, but presented as an explanation of components rather than a timeline.

Elasticsearch
-------------

The Elasticsearch instance (or cluster, depending on the deployment configuration) is used as a database to contain all check definitions, attributes, and check results. Several indices are created in Elasticsearch to store this information.

In the `checkdef` and `check` indices and the `attrib_admin_*` and `attrib_user_*` index patterns, each document represents a single check, and will have a document ID equal to the `id` field in the check's definition.

For example, if a check has an `id` of `dns-active-directory-team01`, then the documents related to that check in the above indices and index patterns will all have document IDs of `dns-active-directory-team01`.

### `checkdef`

This index is the source of check definitions for Dynamicbeat.

### `checks`

This index is a source of generic check information that contains only the `id`, `name`, `type`, and `group` fields of each check definition. This is used for some visualizations that require teams to be able to see what checks are running. This index is useful in addition to the `checkdef` index, because it does not give teams any information about how the checks are configured.

### `attrib_admin_*`

This index pattern contains attributes that can only be read and modified by attribute administrators. Each team has their own admin attribute index. The group name is appended to the index name. For example, a group named `team01` would have an admin attribute index of `attrib_admin_team01`.

### `attrib_user_*`

This index pattern is the same as `attrib_admin_*`, except these attributes can also be read and modified by the team's members. To clarify, only members of `team01` can view and modify the attributes in `attrib_user_team01`.

### `results-admin-*`

This index pattern contains detailed check results for all checks that are running. It makes it easier for competition organizers to search for any check from any team within the Discover app in Kibana and see why a check may have failed.

This index pattern, along with the rest of the `results-*` index patterns, appends a timestamp to the index with the current date in `YYYY.MM.dd` format.

### `results-all-*`

This index pattern also contains the check results for all checks, but the `message` and `details` fields are removed. This lets anybody see what checks are passing for each team (which is required for some dashboards). However, nobody consulting this index pattern would be able to see the details explaining why a check failed.

### `results-TEAM-*`

These index patterns contain detailed check results for a single team's checks, and gives teams a starting point for troubleshooting their failing checks. `TEAM` is just a placeholder for the team's name. For example, `team01`'s results index pattern would be `results-team01-*`.

Kibana
------

The stock Kibana application is what teams use to view the current and past status of their checks, get information about why their checks are failing, and compare themselves to the other teams.

The Kibana instance is configured with two Spaces: Default and Scorestack. The Default space is the normal Kibana default space, and all Kibana applications are enabled in it. The Default space should only be accessible by Scorestack administrators. The Scorestack space only enables the Kibana features that are needed for Scorestack. This is to reduce clutter in the sidebar and prevent confusion for new Scorestack users.

Since dashboards, visualizations, index patterns, and other Kibana saved objects aren't shared across spaces, these saved objects are duplicated across the two spaces. The practical effect of this is that any changes made to an object (like a dashboard) in one space must also be made to the same object in the other space.

Kibana Plugin
-------------

The Kibana plugin adds a new application to Kibana that allows admins and users to view and modify check attributes. The application allows you to view the attributes for your checks, see their current values, and change their values at any time. Once an attribute value has been updated through the plugin application it will immediately update the value in Elasticsearch, but it may take a few rounds of checks for Dynamicbeat to re-read the check definitions and the change to take effect.

Dynamicbeat
-----------

Dynamicbeat is the program that runs check definitions. It built using the libbeat library from Elastic, which is the same library used for other beats like Filebeat, Winlogbeat, and Journalbeat.

Dynamicbeat uses periods to separate "rounds" of checks. A period is a certain amount of time to wait before starting the next round. By default, Dynamicbeat's period is 30 seconds, so every 30 seconds another "round" of checks will be started.

At startup, Dynamicbeat will first query Elasticsearch for all check definitions and check attributes, and then save the results.

Every period, Dynamicbeat will start all checks that it knows about at the same time and run them asynchronously. Once all checks have been completed or a timeout has been hit, whichever happens first, the check results will be created and indexed into Elasticsearch. By default, the global timeout for all checks is 30 seconds. If a check does not finish within those 30 seconds, the check will be automatically marked as failing.

Dynamicbeat performs the following steps to index a check result in Elasticsearch:

1. Convert the boolean `passed` field to an integer in the new `passed_int` field. If a check has passed, `passed_int` will be 1, otherwise it will be 0.
2. Create an `epoch` field that contains the integer value of the check result's `@timestamp` converted to [Epoch time](https://en.wikipedia.org/wiki/Unix_time).
3. Index a copy of the check result in both the `results-admin-*` and `results-GROUP-*` indicies, where `GROUP` is the value of the `group` field in the check definition.
4. Index a copy of the check result with the `message` and `details` fields removed in the `results-all-*` index.

Once Dynamicbeat has started a round, it will re-query Elasticsearch for the latest check definitions and check attributes, and save the results for the next round of checks. If Dynamicbeat has any issues loading the latest check definitions (for example, if Elasticsearch is unreachable), then it will reuse the check information from the previous round.