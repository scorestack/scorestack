Deployment
==========

This section explains how to deploy Scorestack using the preconfigured deployment options. It does not explain exactly how the deployment process works, however there should be plenty of information here for someone new to Scorestack to successfully deploy a new Scorestack instance and get scoring working for their competition.

If have any problems at all deploying Scorestack while following this documentation, please [open a new issue](https://github.com/scorestack/scorestack/issues/new/choose) so we can help you out!

Deploying the Elastic Stack properly for Scorestack is a fairly involved task, so the Scorestack developers maintain automation to deploy Scorestack to a variety of platforms with multiple different possible architectures. An "architecture" is a description of how Scorestack is deployed, and usually just refers to the number of replicas that are deployed for different services. This section will also cover the available architectures, what they look like, and what methods are available for deploying those architectures.