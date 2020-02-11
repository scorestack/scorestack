# scorestack

A security competition scoring system built on the Elastic stack.

## Building dynamicbeat

For a number of reasons, we use vscode dev container to build dynamicbeat. Once
you open the project in the provided dev container, run the following commands
to build the beat:

```shell
cd dynamicbeat
make setup
go get
mage build
```

## Running everything

After dynamicbeat has been installed, run the following commands to set up an
instance of ScoreStack with the example check configurations loaded up and run
an instance of dynamicbeat against those checks.

```shell
docker-compose -f docker/build-certs/docker-compose.yml up
docker-compose up -d
docker/setup.sh
dynamicbeat/dynamicbeat -e -d "*" --path.config docker
```