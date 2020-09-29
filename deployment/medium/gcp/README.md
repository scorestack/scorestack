Medium: GCP
===========

This option will deploy the servers for the Medium architecture on GCP. This deployment is estimated to cost $0.50 per hour - see [this GCP pricing calculator](https://cloud.google.com/products/calculator/#id=0fa433ed-d279-4773-9c7b-358eacffe7d4) for more information. Please note that the "per 1 month" prices listed in the calculator are actually the per-hour costs.

Please note that this deployment will use 32 CPU cores, which is the default global quota for CPU cores on GCP accounts. _You will have to stop or terminate all other instances in your GCP account to use this deployment._

Deploying
---------

First, [set up Terraform with the Google Provider](https://www.terraform.io/docs/providers/google/guides/getting_started.html) and make sure `git` is installed on your system. Then, clone the Scorestack repository to the host you will be deploying Scorestack from.

```shell
git clone https://github.com/scorestack/scorestack.git
cd scorestack/deployment/medium/gcp
```

Next, provide values for the five unset variables within `variables.tf`. Here is a brief description of the variables:

- `project`: The GCP project ID to which Scorestack will be deployed
- `credentials_file`: The path to the GCP credentials file - see the [GCP provider reference](https://www.terraform.io/docs/providers/google/guides/provider_reference.html#credentials) for more information
- `ssh_pub_key_file`, `ssh_priv_key_file`: The paths to the SSH keypair that will be added to the created instances - these must already be created
- `fqdn`: The domain name that Scorestack will be served behind - the DNS record for this domain must be configured manually after deployment

Finally, run Terraform to deploy Scorestack.

```shell
terraform init
terraform apply
```

Once deployment is finished, configure the DNS record for your FQDN to point at the public IP of the `nginx` instance. Then, run the [Ansible playbook](../ansible/README.md).

Destroying
----------

To teardown and destroy the Scorestack cluster, use Terraform.

```shell
terraform destroy
```