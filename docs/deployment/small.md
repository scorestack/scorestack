Small Architecture
==================

This architecture deploys three hosts:

- `elasticsearch`: A single-node Elasticsearch cluster
- `kibana`: The Kibana server
- `logstash`: The Logstash server

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

Once all the containers have started, you can check on the progress of the deployment by following the logs of the setup container. You can do so with another `docker-compose` command.

```shell
sudo docker-compose -f deployment/small/docker/docker-compose.yml logs -f setup
```

If it looks like the setup container is getting hung up on waiting for Elasticsearch or Kibana, give it two to three minutes before troubleshooting. Sometimes it can take a while for these components to get ready. If it's still trying to connect after a few minutes have passed, check the logs for the Elasticsearch and Kibana containers to see what's wrong.

Once the setup container has exited, Scorestack is fully deployed and ready to use!

### Dynamicbeat Certificates

To deploy Dynamicbeat, you will need the Dynamicbeat client certificate and key and the Scorestack internal CA certificate. Once the deployment has finished, these files can be found at the following paths within the repository:

- `deployment/small/docker/certificates/dynamicbeat/dynamicbeat.key`
- `deployment/small/docker/certificates/dynamicbeat/dynamicbeat.crt`
- `deployment/small/docker/certificates/ca/ca.crt`

See the [Dynamicbeat deployment guide](dynamicbeat.md) for more information.

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