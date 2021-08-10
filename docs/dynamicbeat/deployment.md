Deployment
==========

Once Scorestack has been deployed and you've uploaded some checks, Dynamicbeat needs to be deployed to a host for checks to run. Currently, the only supported way of deploying Dynamicbeat is by running it as a systemd service on a Linux system.

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

Service Setup
-------------

The recommended way of running Dynamicbeat is via a systemd service. To configure the service, you can use the following systemd unit file:

```
{{#include dynamicbeat.service}}
```

> This file is also available for download [here](./dynamicbeat.service).

If using the above service file, you should save it to `/etc/systemd/system/dynamicbeat.service`. You can then start Dynamicbeat by running the following commands:

```shell
sudo systemctl daemon-reload
sudo systemctl start dynamicbeat.service
sudo systemctl enable dynamicbeat.service
```

Note that this service file assumes that the Dynamicbeat binary is placed at `/opt/dynamicbeat/dynamicbeat`, and that the configuration file is  at `/opt/dynamicbeat/dynamicbeat.yml`. If your configuration file is placed somewhere else or named differently, add the `--config /path/to/dynamicbeat/config` argument to the `ExecStart` line in your service file.