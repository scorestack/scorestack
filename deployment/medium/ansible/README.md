Medium: Ansible
===============

This option will configure the servers deployed by a Terraform option for the Medium architecture.

Deploying
---------

First, deploy the Medium architecture using one of the available Terraform options. Then, make sure Ansible is installed on your system, and move into this directory.

```shell
cd ${SCORESTACK_PATH}/deployment/medium/ansible
```

Run Ansible to deploy Scorestack.

```shell
ansible-playbook playbook.yml -i inventory.ini
```