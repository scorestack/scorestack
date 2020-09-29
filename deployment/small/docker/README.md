Small: Docker
=============

This option will deploy the Small architecture on a single host using Docker.

Deploying
---------

First, make sure `docker`, `docker-compose`, and `git` are installed on your system. Then, clone the Scorestack repository to the host you plan to use for running Scorestack.

```shell
git clone https://github.com/scorestack/scorestack.git
cd scorestack
```

Next, you need to increase the maximum number of memory map areas that a single process can use. This is required for Elasticsearch to function properly. If you forget to configure this setting, Elasticsearch will appear to start properly, but crash and exit before the APIs are available. This setting is configured via `sysctl` on Linux hosts.

You can change this setting immediately with a `sysctl` command.

```shell
sudo sysctl -w vm.max_map_count=262144
```

Please note the changes made by the above command will not persist on reboot. In order to persistently make these changes, add the `vm.max_map_count` setting to your `/etc/sysctl.conf` file.

Finally, run docker-compose to deploy Scorestack.

```shell
sudo docker-compose -f deployment/small/docker/docker-compose.yml up -d
```

Destroying
----------

To teardown and destroy the Scorestack cluster, use docker-compose.

```shell
sudo docker-compose -f deployment/small/docker/docker-compose.yml down -v
```

If you would like to remove the generated certificates, remove the certificates folder.

```shell
sudo rm -rf deployment/small/docker/certificates
```