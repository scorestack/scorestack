Scorestack
==========

![Dynamicbeat](https://github.com/scorestack/scorestack/workflows/Dynamicbeat/badge.svg)
![Kibana Plugin](https://github.com/scorestack/scorestack/workflows/Kibana%20Plugin/badge.svg)

A security competition scoring system built on the Elastic stack.

Quickstart
----------

To get up and running checking your services with Scorestack as fast as possible, follow the deployment guide for the [small deployment with Docker](./deployment/small.md#docker-deployment). After that, read up on [adding checks](./checks/adding_checks.md), and then take a look at the example checks in the [`examples` folder in the repository](https://github.com/scorestack/scorestack/tree/stable/examples) - you should be able to find examples to get you started for several of your services! After you've added some checks, check out the [Dynamicbeat deployment guide](./deployment/dynamicbeat.md) to get Dynamicbeat up and running.

Documentation Overview
----------------------

This documentation is split into sections based on overall topics and audiences. Each section starts with a brief introduction that explains in more detail its intended audience and what it covers.

### Checks

The [check guide](./checks.md) and [check reference](./reference.md) cover pretty much everything you need to know about writing and adding checks to Scorestack. These sections are written for anyone who is writing checks for Scorestack - not just people who are managing a Scorestack instance.

### User Guides

The [deployment guide](./deployment.md) covers everything you need to know about deploying a Scorestack instance. It also explains the various automation and architectures that are maintained by the Scorestack developers.

### Development Documentation

The development documentation consists of the [design docs](./design.md) and the [binary building guide](./building.md). These sections are primarily written for Scorestack contributors, but Scorestack administrators may also find them useful. The binary building guide explains how to build your own copies of the binaries provided in the Scorestack releases. The design docs attempt to comprehensively describe how Scorestack works.

### Appendices

The appendices are the catch-all for any other content that doesn't fall into any of the above categories. For now this includes an [FAQ](./appendix/faq.md), but other topics will be added over time.