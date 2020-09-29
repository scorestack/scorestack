Medium
======

The Medium architecture deploys seven hosts:

- `elasticsearch1`, `elasticsearch2`, `elasticsearch3`: Elasticsearch master-eligible data nodes
- `elasticsearch4`: An Elasticsearch coordinating-only node
- `kibana`: The Kibana server
- `logstash`: The Logstash server
- `nginx`: An Nginx server that proxies public requests to Kibana, Logstash, and the Elasticsearch cluster. Also used as an SSH jump box.

Deployments
-----------

The Medium architecture uses a two-stage deployment mechanism. First, Terraform is used to deploy servers on a cloud or virtualization provider. Then, Ansible is used to configure those servers.

The following Terraform options are available.

- [GCP](./gcp/README.md)

After running the Terraform of your choice, run [the Ansible playbook](./ansible/README.md) to configure the servers.