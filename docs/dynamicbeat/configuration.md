Configuration
=============

Dynamicbeat pulls configuration values from three sources: environment variables, command-line arguments, and configuration files. These sources are combined in the following decending priority order to create a single configuration:

1. Command-line arguments
2. Environment variables
3. Configuration file
4. Default values

A setting configured via a command-line argument will override all other settings.

Creating a Config File
----------------------

While Dynamicbeat doesn't require a configuration file, it is usually desirable to store the Dynamicbeat configuration in a configuration file. To save the current Dynamicbeat configuration to a file named `dynamicbeat.yml`, run:

```shell
dynamicbeat config save dynamicbeat.yml
```

This will parse any command-line arguments or environment variables you've set to configure Dynamicbeat (if any), and then render a YAML file containing the configuration. Feel free to edit the configuration file as necessary.

Alternatively, you can use the [reference configuration file](#configuration-reference) as a starting point.

Viewing Your Config
-------------------

If you would like to view your current configuration (and don't want to save it to a file), you can use the `confg view` subcommand:

```shell
dynamicbeat config view
```

Your current configuration will be printed in YAML format.

Specifying a Config File
------------------------

Dynamicbeat will automatically use a configuration file in the current directory named `dynamicbeat.yml` (or files with other supported extensions - see [file formats](#other-config-file-formats)). You can also tell Dynamicbeat which config file to use with the `--config` flag:

```shell
dynamicbeat run --config /path/to/dynamicbeat/config
```

Other Config File Formats
-------------------------

If you don't like YAML, Dynamicbeat supports other config formats as well! When running `dynamicbeat config save`, try passing a filepath with other extensions like `.json`, `.toml`, or `.env`. If you use an extension that Dynamicbeat doesn't support, it will print a message listing the currently supported config file types.

For example, to save the config in TOML format, run the following command:

```shell
dynamicbeat config save dynamicbeat.toml
```

Configuration Reference
-----------------------

This reference can also be downloaded [here](dynamicbeat.yml) as a YAML file.

```yaml
{{#include dynamicbeat.yml}}
```