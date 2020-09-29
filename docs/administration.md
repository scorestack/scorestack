Scorestack Administration
=========================

This document explains how to deploy and manage your Scorestack environment. Please note that this document does not cover writing check definitions. Instead, this document focuses on how to obtain the necessary binaries for running Scorestack, deploying an Elastic Stack cluster that will support Scorestack, and running Dynamicbeat against the deployed Elastic Stack.

Architectural Overview
----------------------

Scorestack is based around a customized deployment of the Elastic Stack that includes Elasticsearch, Kibana, Logstash, a Kibana plugin, and a custom Beat, named Dynamicbeat. All of these components are configured with X-Pack security and TLS client/server verification is used for all inter-cluster communications.

### Running Checks

A **check** is a single attempt to verify the functionality of a specific network service. Here are a few examples of basic checks that could be run:

- Send a series of ICMP Echo Requests to an IP address and expect a certain number or percentage of ICMP Echo Replies from that address.
- Send an HTTP GET request for a specific webpage to a webserver and check that the returned content matches the expected content.
- Log in to a system via SSH with a given set of credentials, run a specific command, and ensure the command printed the expected content.

Each check is defined by a **check definition**, which is a JSON document stored in Elasticsearch that provides the necessary information for Dynamicbeat to run the check. Additionally, a check may have **check attributes**, which allow for on-the-fly modification of variables that are templated into the check definition. For more information on defining checks, see [the check definition documentation](./checks.md).

When Dynamicbeat first starts, it pulls the check definitions stored in Elasticsearch and stores them in memory. Then every 30 seconds, Dynamicbeat will start a single check for each one of the check definitions that it currently has stored in memory. Additionally, every minute Dynamicbeat will refresh the check definition information that it has stored in memory by querying Elasticsearch for updates.

### Reporting Check Results

As checks finish executing, their results (pass, fail, or timeout, with some additional information) are buffered to be sent to Logstash later on. Typically, Dynamicbeat will immediately send check results to Logstash. However, if there issues establishing a stable connection to Logstash, the results will stay in Dynamicbeat's buffer until it can reestablish the connection.

Once Logstash receives a check result from Dynamicbeat, it will perform some basic processing on the event. First, fields are pruned from the event so that only Scorestack-specific fields are included; this removes a variety of metadata fields that the Beats framework adds to events, but aren't relevant to Scorestack.

Then, the `passed` boolean field is converted to an integer in the `passed_int` field. If the check passed, `passed_int` will be set to `1`. Otherwise, `passed_int` will be `0`. This conversion allows for easy score calculation within Kibana dashboards.

Next, the `@timestamp` field is converted to an integer representing the Unix epoch representation of the timestamp, which is stored in the `epoch` field. This conversion makes it simple to display only the latest check results within Kibana dashboards.

Finally, three versions of the result event are created: generic, admin, and group. These events are then stored in an Elasticsearch index that matches the glob `results-*-TIMESTAMP`, where `TIMESTAMP` is a timestamp representing the current date in the format `YYYY.MM.DD`.

#### Generic Results

Generic results have the `message` and `details` fields removed, and are viewable by all Scorestack users. This allows teams to see how other teams are doing, but does not give them information on _why_ other teams' checks may be failing. Since field-based access control is a premium feature of the Elastic Stack, this workaround is required for competition-wide dashboards to work without revealing details of check results to other teams.

Generic results are stored in the `results-all-*` indices.

#### Group Results

The group results do not have the `message` and `details` fields removed, and are only viewable by members of the group the check belongs to. For example, the group results for a check in the Team 3 group can only be viewed by users who are in the Team 3 group. Group results allow users to get more detailed information about why their checks may be failing, which provides them a starting point for troubleshooting their services.

Group results are stored in the `results-GROUP-*` indices, where `GROUP` is the group ID.

#### Admin Results

The admin results, like the group results, do not have the `message` and `details` fields removed. However, admin results are only viewable by members of the `spectator` group. Admin results are mainly useful for troubleshooting issues with service deployment prior to a competition, or for detecting issues with check definitions and/or Dynamicbeat. This is becase all the admin results are stored within a single set of indices, unlike the group events, which are stored across a variety of indices. Having all admin results in a single set of indices allows Scorestack administrators to search across all check results using a single index glob.

The admin results are stored in the `results-admin-*` indices.

### Certificates

All inter-cluster communications are secured using a private TLS certificate chain that is typically generated during deployment. All clients and servers that are private to the cluster are required to present valid TLS certificates that are signed by the cluster root CA. The only communications which are not secured in this manner are those to the Elasticsearch HTTPS API, since the API uses username/password authentication.

While most inter-cluster communications should be firewalled off from any traffic external to the cluster, the Logstash input for Dynamicbeat's check results must be publicly facing, so this is the most important endpoint to require TLS authentication. If the Logstash input is misconfigured to accept client certificates that aren't signed by the cluster root CA, then malicious users can submit fraudulent check results.

Obtaining Binaries
------------------

In order to run Dynamicbeat and Scorestack's Kibana plugin, the binaries for these components must be obtained. The Dynamicbeat binary is a compiled Golang executable that runs Dynamicbeat. The Kibana plugin binary is a zipfile containing the compiled assets that are installed into the Kibana server and loaded at runtime.

### Prebuilt Binaries

Most users will want to use the prebuilt binaries that are available on the [Scorestack Releases page](https://github.com/scorestack/scorestack/releases). The Kibana plugin zipfile and a zipped Dynamicbeat executable are attached to each release. This is the recommended way of obtaining the binaries.

### Building Your Own Binaries

If you really want to, you can build these binaries yourself. Please see the [documentation on building](./building.md) for more information on how to do this.

Deploying a Cluster
-------------------

Since deploying the Elastic Stack properly for Scorestack is a fairly involved task, the Scorestack developers maintain automation to deploy Scorestack to a variety of platforms and architectures. Please see the [deployment documentation](../deployment/README.md) for more information on using this automation.

Deploying Scorestack manually is **not** recommended or supported. Please use the provided automation to deploy Scorestack. If the provided automation does not meet your needs, please submit a new issue and explain what automation you would like to be added.

Running Dynamicbeat
-------------------

Once Scorestack has been deployed and you have a Dynamicbeat binary, Dynamicbeat must be deployed to a host properly to run checks and communciate with Scorestack. Currently, the only supported way of deploying Dynamicbeat is by running it as a systemd service on a Linux system.

### System Requirements

Currently, Dynamicbeat is only supported on Linux systems. Please note that Dynamicbeat can run checks _against_ non-Linux systems (like Windows or Mac), but Dynamicbeat itself must run on Linux.

The system that Dynamicbeat runs on must be able to access the Elasticsearch HTTP API and the Logstash Beat input endpoints. Additionally, the system should be able to access the services that Dynamicbeat will be running checks against. Varying the placement of Dynamicbeat on the network can help mask the check traffic slightly. For example, running Dynamicbeat on the same network that red team traffic originates from is a common practice.

Dynamicbeat can typically run fine with minimal system resources. In the past, Dynamicbeat has worked fine with only 2 CPU cores and 2 GB of RAM, but you should always perform some testing to determine what will work best for your environment.

#### ICMP Permissions

The ICMP library that Dynamicbeat's ICMP protocol uses requires specific sysctl parameters to be set in order to work. The following command can be used to immediately set the sysctl parameter, but the change will not persist on a reboot. You should set the parameter in `/etc/sysctl.d` (or the equivalent on your distribution) in order for changes to persist on reboot.

```shell
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
```

For more information on the ICMP permissions required, see [the library's documentation](https://github.com/go-ping/ping#note-on-linux-support).

#### Dynamicbeat Disk Impact

When the systems that Dynamicbeat runs checks against are using shared network storage, the available IOPS is very important, especially for larger competitions. This is because each check that Dynamicbeat runs can trigger multiple disk I/O operations on the system being checked. Additionally, checks are executed nearly simultaneously, so a round of checks can have a serious impact on disk performance.

If possible, try stress testing your infrastructure before the start of the competition by running multiple instances of Dynamicbeat simultaneously. While running this stress test, try interacting with the competition systems. If your competition has any "Desktop Experience" (GUI) Windows systems, try interacting with them - opening multiple programs, performing searches in the Start menu, etc. - since the experience on Windows systems can be particularly sensitive to a lack of sufficient IOPS.

If you do experience IOPS issues, one possible solution is to increase the time between rounds. By default rounds run every 30 seconds, but increasing this to 2 or 3 minutes (or even more) can help.

### Configuring Dynamicbeat

Dynamicbeat is configured via a YAML configuration file that is largely similar to configuration files for other Beats, like Filebeat and Metricbeat. The following section provides an overview of the most common configuration options:

```yaml
dynamicbeat:
  # Defines how often checks will be executed
  #period: 30s

  # Defines how often check definitions and attributes will be read from elasticsearch
  #update_period: 1m

  # Where check definitions will be read from
  #check_source:
    # List of Elasticsearch hosts to query
    #hosts: ["https://localhost:9200"]
    
    # Credentials used by this beat to read check definitions
    #username: dynamicbeat
    #password: changeme

    # Whether to verify SSL/TLS certificates, if HTTPS is used
    #verify_certs: true

    # The index that check definitions will be read from
    #index: checkdef

output.logstash:
  # The Logstash hosts
  #hosts: ["localhost:5454"]

  # Enable SSL.
  ssl.enabled: true

  # List of root certificates for HTTPS server verifications
  #ssl.certificate_authorities: ["/etc/pki/root/ca.pem"]

  # Certificate for SSL client authentication
  #ssl.certificate: "/etc/pki/client/cert.pem"

  # Client Certificate Key
  #ssl.key: "/etc/pki/client/cert.key"
```

Please note that the port for `output.logstash.hosts` _should_ be set to `5454`, not the default of `5044`. By default, the Dynamicbeat pipeline for Logstash listens on port 5454.

Here is an example Dynamicbeat configuration:

```yaml
dynamicbeat:
  check_source:
    hosts: ["https://demo.scorestack.io:9200"]
    username: dynamicbeat
    password: example
output.logstash:
  hosts: ["demo.scorestack.io:5454"]
  ssl.enabled: true
  ssl.certificate_authorities: ["/opt/dynamicbeat/ca.pem"]
  ssl.certificate: "/opt/dynamicbeat/cert.pem"
  ssl.key: "/opt/dynamicbeat/cert.key"
```

### Service Setup

The recommended way of running Dynamicbeat is via a systemd service. To configure the service, you can use the following systemd unit file:

```
[Unit]
Description=Dynamicbeat
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/dynamicbeat
ExecStart=/opt/dynamicbeat/dynamicbeat

[Install]
WantedBy=multi-user.target
```

If using the above service file, you should save it to `/etc/systemd/system/dynamicbeat.service`. You can then start Dynamicbeat by running the following commands:

```shell
sudo systemctl daemon-reload
sudo systemctl start dynamicbeat.service
sudo systemctl enable dynamicbeat.service
```

Note that this service file assumes that the Dynamicbeat binary is placed at `/opt/dynamicbeat/dynamicbeat`, and that the configuration file is placed at `/opt/dynamicbeat/dynamicbeat.yml`. The key and certificates for Dynamicbeat should be stored at the locations referenced within the `output.logstash.ssl.*` options in the Dynamicbeat configuration file.

User Permissions
----------------

This section describes the roles that are added to Scorestack, their permissions, and their intended uses.

### `dynamicbeat_reader`

This role provides read-only access to the `checkdef*` and `attrib_*` indices. This role is intended to be used by the Dynamicbeat user, and provides Dynamicbeat with the bare minimum amount of permissions required for proper operation.

### `logstash_writer`

This role provides access to the `results-*` indices and some cluster permissions that are required by Logstash to connect properly to Elasticsearch. This role provides the bare minimum amount of permissions required for Logstash, and should only be used by the Logstash user.

### `common`

This role provides read-only access to the `results-all*` and `checks` indices, and provides read access to the `scorestack` space. This allows users to view the generic results and some generic checks, which is required for the overall dashboards to work properly. The `scorestack` space is a customized Kibana space that only includes Scorestack-specific components, reducing the clutter in the Kibana UI. Read-only access is provided to only that space so that users don't have to pick between it and the default space, which has several components included that are not needed for Scorestack.

This role should be used for all Scorestack end-users that interact with Kibana.

### `spectator`

This role provides read-only access to the `results*` indices. This allows users to view team-specific dashboards and admin/group check results. This role should generally only be given to Scorestack administrators or "spectator" users like redteam and whiteteam.

### `attribute-admin`

This role provides full access to the `attrib_*` indices. This allows users to modify all attributes of all teams. This role should only be given to Scorestack administrators that need to modify administrator attributes, or assist teams with modifying their own attributes.

### `check-admin`

This role provides full access to the `check*` indicies. This allows users to create, modify, and delete all check definitions. This role should only be given to Scorestack administrators that are managing checks.

### Group Roles

A role is created for each group that gets added, which provides read access to the group results index for the group and read/write access to the group's user attributes index. This allows group users to see detailed check results for their group, view team-specific dashboards for their group, and modify user attributes for their group. This role should only be given to the associated group user.

Notes
-----

This section explains things you may encounter while running Scorestack that don't really fit anywhere else.

### Shard Failures

When using a multi-node Elasticsearch cluster (like with the [Medium deployment](../deployment/medium/README.md)), the Elasticsearch nodes must be able to communicate with each other without any issues. If there are any network connectivity problems between the nodes such as packet loss, the Elasticsearch nodes may lose contact with each other and think all shards have failed. Usually after a little while, the Elasticsearch nodes will be able to resolve the issue themselves, but this will cause Scorestack to briefly become unavailable.