Dynamicbeat
===========

This section of the documentation covers the usage, configuration, and deployment of Dynamicbeat.

The primary function of Dynamicbeat is to run the checks that you've defined in your Scorestack instance, but it also serves as an administration utility for Scorestack instances.

The [Configuration page](./dynamicbeat/configuration.md) has more information about how to configure Dynamicbeat, and contains a heavily-commented YAML configuration reference.

The [Deployment page](./dynamicbeat/deployment.md) provides an overview of the options to deploy a Dynamicbeat instance for running your checks.

The [Overrides page](./dynamicbeat/overrides.md) explains the purpose of the [team overrides system](./checks/file_format.md#team-overrides). A few examples are provided to help you write check definitions that work across _all_ your teams.

Finally, the [Reference section](./dynamicbeat/reference/dynamicbeat.md) contains automatically-generated documentation covering the command-line interface of Dynamicbeat. This information is also available via the `-h` and `--help` flags when running Dynamicbeat.