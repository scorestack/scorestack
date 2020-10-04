Building Binaries
=================

Scorestack relies on two different binaries: a Golang executable for Dynamicbeat, and a zipfile for the Kibana plugin. While prebuilt versions of these binaries are attached to each release, this document provides an explanation of how to build your own binaries if you so choose.

Dynamicbeat
-----------

To build and test Dynamicbeat with the provided Docker container, run the following command:

```shell
docker-compose -f dockerfiles/docker-compose.yml run dynamicbeat-ci /scripts/test.sh
```

The compiled binary can be found at `dynamicbeat/dynamicbeat`.

### Building During Development

If you are actively developing Dynamicbeat, you may wish to build Dynamicbeat within the devcontainer to test your changes. Instead of using the above docker-compose command, you can perform the following steps within the `dynamicbeat/` folder.

First, install the required go dependencies:

```shell
go get
```

Then, build the binary:

```shell
make
```

The compiled binary will be named `dynamicbeat` within the `dynamicbeat/` folder.

Kibana Plugin
-------------

To build the Kibana plugin with the provided Docker container, run the following command:

```shell
docker-compose -f dockerfiles/docker-compose.yml run kibana-plugin-ci /scripts/test.sh
```

Please note that this build process can take quite a long time (20-30 minutes).

The zipfile can be found at `kibana-plugin/build/kibana-plugin-0.0.0.zip`.

### Building During Development

Typically during Kibana plugin development, you will want to use the Kibana development server to test the plugin. However, sometimes it can be useful to build the zipfile during development to ensure it can properly install and run in a production version of Kibana. If working within the Kibana plugin devcontainer, you run the following commands in the `kibana-plugin/` folder.

First, make sure you're in the right directory:

```shell
cd $HOME/kibana/plugins/scorestack
```

Next, install all the required dependencies.

```shell
yarn kbn bootstrap
```

The first time you run `yarn kbn bootstrap`, it can take several (20+) minutes to install all the required dependencies. However, after the first time you run it, running the same command should be fairly quick later on.

With the dependencies installed, you can build the zipfile.

```shell
yarn plugin-helpers build
```

The zipfile can be found at `kibana-plugin/build/kibana-plugin-0.0.0.zip`.