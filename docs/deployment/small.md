Small Architecture
==================

This architecture deploys three hosts:

- `elasticsearch`: A single-node Elasticsearch cluster
- `kibana`: The Kibana server

Docker Deployment
-----------------

The Small architecture can deployed on a single system using Docker. Before deploying with Docker, you must make sure that `docker`, `docker-compose`, and `git` are installed on your system. Then, [clone the Scorestack repository](cloning.md) to the host you plan to use for running Scorestack.

### Deploying

First, change into the Scorestack repository directory.

```shell
cd scorestack
```

Next, you need to increase the maximum number of memory map areas that a single process can use. This is required for Elasticsearch to function properly. If you forget to configure this setting, Elasticsearch will appear to start properly, but crash and exit before the APIs are available. This setting is configured via `sysctl` on Linux hosts.

You can change this setting immediately with a `sysctl` command.

```shell
sudo sysctl -w vm.max_map_count=262144
```

Please note the changes made by the above command will not persist on reboot. In order to persistently make these changes, add the `vm.max_map_count` setting to your `/etc/sysctl.conf` file, or to a file under the `/etc/sysctl.d` directory. The proper location for these changes depends on your distribution.

Now run `docker-compose` to deploy Scorestack. Make sure to include the `-d` parameter!

```shell
sudo docker-compose -f deployment/small/docker/docker-compose.yml up -d
```

Then, configure Scorestack settings in the Elastic Stack using [Dynamicbeat's `setup` command](./setup.md).

### Cleanup/Teardown

If you just want to stop Scorestack and restart it later, you can do so with normal `docker-compose` commands.

```shell
sudo docker-compose -f deployment/small/docker/docker-compose.yml down
sudo docker-compose -f deployment/small/docker/docker-compose.yml up -d
```

> If you restart Scorestack, the `setup` container will also be restarted, and that's okay. It shouldn't make any changes to your cluster.

To destroy the Scorestack cluster completely and remove all artifacts, you must stop and remove the containers, remove the volumes, and delete the certificates from disk.

```shell
sudo docker-compose -f deployment/small/docker/docker-compose.yml down -v
rm -rf deployment/small/docker/certificates
```