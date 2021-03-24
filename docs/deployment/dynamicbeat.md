Dynamicbeat
===========

Once Scorestack has been deployed and you have a Dynamicbeat binary, Dynamicbeat must be deployed to a host for checks to be run. Currently, the only supported way of deploying Dynamicbeat is by running it as a systemd service on a Linux system.

System Requirements
-------------------

Currently, Dynamicbeat is only supported on Linux systems. Please note that Dynamicbeat can run checks _against_ non-Linux systems (like Windows or Mac), but Dynamicbeat itself must run on Linux.

The system that Dynamicbeat runs on must be able to access the Elasticsearch HTTP API endpoints. Additionally, the system should be able to access the services that Dynamicbeat will be running checks against. Varying the placement of Dynamicbeat on the network can help mask the check traffic slightly. For example, running Dynamicbeat on the same network that red team traffic originates from is a common practice.

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

Dynamicbeat pulls configuration values from three sources: environment variables, command-line arguments, and configuration files. These sources are combined in the following decending priority order to create a single configuration:

1. Command-line arguments
2. Environment variables
3. Configuration file
4. Default values

A setting configured via a command-line argument will override all other settings.

### Creating a Config File

While Dynamicbeat doesn't require a configuration file, it is usually desirable to store the Dynamicbeat configuration in a configuration file. To save the current Dynamicbeat configuration to a file named `dynamicbeat.yml`, run:

```shell
dynamicbeat config save dynamicbeat.yml
```

This will parse any command-line arguments or environment variables you've set to configure Dynamicbeat (if any), and then render a YAML file containing the configuration. Feel free to edit the configuration file as necessary.

The default YAML configuration looks like this:

```yaml
elasticsearch: https://localhost:9200
password: changeme
round_time: 30s
username: dynamicbeat
verify_certs: false
```

> This configuration will work if you run Dynamicbeat on the same host as a default small/docker deployment.

### Viewing Your Config

If you would like to view your current configuration (and don't want to save it to a file), you can use the `confg view` subcommand:

```shell
dynamicbeat config view
```

Your current configuration will be printed in YAML format.

### Specifying a Config File

Dynamicbeat will automatically use a configuration file in the current directory named `dynamicbeat.yml` (or files with other supported extensions - see [file formats](#other-config-file-formats)). You can also tell Dynamicbeat which config file to use with the `--config` flag:

```shell
dynamicbeat run --config /path/to/dynamicbeat/config
```

### Other Config File Formats

If you don't like YAML, Dynamicbeat supports other config formats as well! When running `dynamicbeat config save`, try passing a filepath with other extensions like `.json`, `.toml`, or `.env`. If you use an extension that Dynamicbeat doesn't support, it will print a message listing the currently supported config file types.

For example, to save the config in TOML format, run the following command:

```shell
dynamicbeat config save dynamicbeat.toml
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
ExecStart=/opt/dynamicbeat/dynamicbeat run

[Install]
WantedBy=multi-user.target
```

If using the above service file, you should save it to `/etc/systemd/system/dynamicbeat.service`. You can then start Dynamicbeat by running the following commands:

```shell
sudo systemctl daemon-reload
sudo systemctl start dynamicbeat.service
sudo systemctl enable dynamicbeat.service
```

Note that this service file assumes that the Dynamicbeat binary is placed at `/opt/dynamicbeat/dynamicbeat`, and that the configuration file is placed at `/opt/dynamicbeat/dynamicbeat.yml`.

Settings
--------

This section lists all configurable Dynamicbeat settings in alphabetical order.

### Elasticsearch

Address of the Elasticsearch host to pull checks from and store results in.

- **Type**: String
- **Configuration Key**: `elasticsearch`
- **Environment Variable**: `ELASTICSEARCH`
- **Command-Line Argument**: `-e / --elasticsearch`
- **Default**: `https://localhost:9200`

### Password

Password to use when authenticating with Elasticsearch.

- **Type**: String
- **Configuration Key**: `password`
- **Environment Variable**: `PASSWORD`
- **Command-Line Argument**: `-p / --password`
- **Default**: `changeme`

### Round Time

Time to wait between rounds of checks.

- **Type**: String (must be parsable by [time.ParseDuration](https://golang.org/pkg/time/#ParseDuration))
- **Configuration Key**: `round_time`
- **Environment Variable**: `ROUND_TIME`
- **Command-Line Argument**: `-r / --round_time`
- **Default**: `30s`

### Username

Username to use when authenticating with Elasticsearch.

- **Type**: String
- **Configuration Key**: `username`
- **Environment Variable**: `USERNAME`
- **Command-Line Argument**: `-u / --username`
- **Default**: `dynamicbeat`

### Verify Certificates

Whether to verify the Elasticsearch TLS certificates.

- **Type**: Boolean
- **Configuration Key**: `verify_certs`
- **Environment Variable**: `VERIFY_CERTS`
- **Command-Line Argument**: `-v / --verify_certs`
- **Default**: `false`