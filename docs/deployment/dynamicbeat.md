Dynamicbeat
===========

Once Scorestack has been deployed and you have a Dynamicbeat binary, Dynamicbeat must be deployed to a host for checks to be run. Currently, the only supported way of deploying Dynamicbeat is by running it as a systemd service on a Linux system.

System Requirements
-------------------

Currently, Dynamicbeat is only supported on Linux systems. Please note that Dynamicbeat can run checks _against_ non-Linux systems (like Windows or Mac), but Dynamicbeat itself must run on Linux.

The system that Dynamicbeat runs on must be able to access the Elasticsearch HTTP API and the Logstash Beat input endpoints. Additionally, the system should be able to access the services that Dynamicbeat will be running checks against. Varying the placement of Dynamicbeat on the network can help mask the check traffic slightly. For example, running Dynamicbeat on the same network that red team traffic originates from is a common practice.

Dynamicbeat can typically run fine with minimal system resources. In the past, Dynamicbeat has worked fine with only 2 CPU cores and 2 GB of RAM, but you should always perform some testing to determine what will work best for your environment.

### ICMP Permissions

The ICMP library that Dynamicbeat's ICMP protocol uses requires specific sysctl parameters to be set in order to work. The following command can be used to immediately set the sysctl parameter, but the change will not persist on a reboot. You should set the parameter in `/etc/sysctl.d` (or the equivalent on your distribution) in order for changes to persist on reboot.

```shell
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
```

For more information on the ICMP permissions required, see [the library's documentation](https://github.com/go-ping/ping#note-on-linux-support).

### Dynamicbeat Disk Impact

When the systems that Dynamicbeat runs checks against are using shared network storage, the available IOPS is very important, especially for larger competitions. This is because each check that Dynamicbeat runs can trigger multiple disk I/O operations on the system being checked. Additionally, checks are executed nearly simultaneously, so a round of checks can have a serious impact on disk performance.

If possible, try stress testing your infrastructure before the start of the competition by running multiple instances of Dynamicbeat simultaneously. While running this stress test, try interacting with the competition systems. If your competition has any "Desktop Experience" (GUI) Windows systems, try interacting with them - opening multiple programs, performing searches in the Start menu, etc. - since the experience on Windows systems can be particularly sensitive to a lack of sufficient IOPS.

If you do experience IOPS issues, one possible solution is to increase the time between rounds. By default rounds run every 30 seconds, but increasing this to 2 or 3 minutes (or even more) can help.

Configuring Dynamicbeat
-----------------------

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

Service Setup
-------------

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