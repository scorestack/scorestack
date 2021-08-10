Medium Architecture
===================

This architecture deploys seven hosts:

- `elasticsearch1`, `elasticsearch2`, `elasticsearch3`: Elasticsearch master-eligible data nodes
- `elasticsearch4`: An Elasticsearch coordinating-only node
- `kibana`: The Kibana server
- `nginx`: An nginx server that proxies public requests to Kibana and the Elasticsearch cluster. Also used as an SSH jump box.

GCP Deployment
--------------

The Medium architecture can be deployed to GCP using multiple VMs using Terraform and Ansible. Before running the Terraform for the GCP deployment, make sure to [set up Terraform with the Google Provider](https://www.terraform.io/docs/providers/google/guides/getting_started.html) and make sure `git` is installed on your system. Then, [clone the Scorestack repository](cloning.md) to the host you plan to use for deploying Scorestack. This can be your laptop or wherever you're most comfortable - Scorestack itself will not be running in this host!

> Please note that new GCP accounts and projects have a per-region quota on CPU usage. This deployment will consume the entire quota.

### Deploying

First, change into the GCP terraform directory.

```shell
cd scorestack/deployment/medium/gcp
```

Next, in a new file named `terraform.tfvars`, provide values for the four unset variables defined in `variables.tf`. Here is a brief description of the required variables:

- `project`: The GCP project ID to which Scorestack will be deployed
- `credentials_file`: The path to the GCP credentials file - see the [GCP provider reference](https://www.terraform.io/docs/providers/google/guides/provider_reference.html#credentials) for more information
- `ssh_pub_key_file`, `ssh_priv_key_file`: The paths to the SSH keypair that will be added to the created instances - these must already be created

Here's an example of what your `terraform.tfvars` file might look like:

```ini
project = "scorestack-300023"
credentials_file = "~/.config/gcloud/application_default_credentials.json"
ssh_pub_key_file = "~/.ssh/scorestack.pub"
ssh_priv_key_file = "~/.ssh/scorestack"
```

> If you have a domain that you'd like to point at Scorestack, make sure to set the `fqdn` variable to that domain. For example, if you were to run a Scorestack instance at `demo.scorestack.io`, you would add the following line to your `terraform.tfvars` file:
> ```ini
> fqdn = "demo.scorestack.io"
> ```

Next, install the necessary Terraform components to work with the GCP APIs.

```shell
terraform init
```

If you are using a brand-new GCP project, you will first have to [enable billing](https://cloud.google.com/billing/docs/how-to/modify-project) for the project. Once billing is enabled, you can enable the Compute Engine API, which is required for the Terraform to work properly.

```shell
gcloud services enable compute.googleapis.com
```

Finally, run Terraform to deploy Scorestack.

```shell
terraform apply
```

Once deployment is finished, configure the DNS record for your FQDN to point at the public IP of the `nginx` instance. Then, run the [ansible deployment](#ansible-deployment) for the medium architecture.

### Cleanup/Teardown

To destroy the Scorestack cluster completely and remove all artifacts, you must destroy the Terraform resources, the generated certificates, and the Ansible inventory file.

```shell
terraform destroy
rm -rf ../ansible/certificates
rm ../ansible/inventory.ini
```

Ansible Deployment
------------------

Once the infrastructure for the medium architecture has been deployed by one of the available Terraform options, Ansible is used to deploy and configure Scorestack. You will need to install Ansible on your system for this deployment.

### Deploying

Once your Terraform deployment is complete, change from your Terraform directory to the Ansible directory in the repository.

```shell
cd ../ansible
```

Then, run the Ansible playbook to deploy Scorestack. The inventory file was generated for you by Terraform to properly configure SSH access for Ansible to all of the instances.

```shell
ansible-playbook playbook.yml -i inventory.ini
```

Finally, configure Scorestack settings in the Elastic Stack using [Dynamicbeat's `setup` command](./setup.md).