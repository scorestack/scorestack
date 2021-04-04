Setup
=====

Once your deployment has finished, you will have an almost-stock configuration of Elastic Stack ready to go. In order to configure the users, dashboards, indices, and other settings for Scorestack, you can use Dynamicbeat.

You can download a precompiled version of Dynamicbeat for your operating system from the [latest release on GitHub](https://github.com/scorestack/scorestack/releases/latest).

Save the zipfile and extract Dynamicbeat to the system that you'll be writing check configurations on. Next, you'll need to [configure Dynamicbeat](../dynamicbeat/configuration.md) so it knows how to configure your Scorestack instance.

Running Setup
-------------

To set up everything on both Kibana and Elasticsearch, run Dynamicbeat's `setup` command:

```shell
dynamicbeat setup
```

Kibana will be configured first, followed by Elasticsearch. If Kibana or Elasticsearch aren't ready yet, Dynamicbeat will wait until they're up and running and Dynamicbeat can authenticate to them. If you just ran your Scorestack deployment, this might take a little while.

> If Dynamicbeat has been waiting for Kibana or Elasticsearch for three minutes or longer, check the Kibana and/or Elasticsearch logs for errors. Also, double-check the Elasticsearch URL, Kibana URL, and setup credentials in your [Dynamicbeat configuration file](../dynamicbeat/configuration.md#configuration-reference).

Alternatively, you can configure Kibana and Elasticsearch one-by-one:

```shell
dynamicbeat setup kibana
dynamicbeat setup elasticsearch
```

Make sure to run the Kibana setup first - the Elasticsearch setup will fail if Kibana has not already been configured.

Re-Running Setup
----------------

The setup command is safe to be rerun at any time, as many times as desired. All dashboards and roles will be re-created. If any spaces, indices, or users that were configured still exist, they will be left untouched.

> If a user's password has been changed, their password will not be reset when you re-run the setup command.