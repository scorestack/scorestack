# scorestack

![Dynamicbeat CI](https://github.com/s-newman/scorestack/workflows/Dynamicbeat%20CI/badge.svg)

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

After dynamicbeat has been built, run the following commands to set up an
instance of ScoreStack with the example check configurations loaded up and run
an instance of dynamicbeat against those checks.

```shell
sudo sysctl -w vm.max_map_count=262144
docker-compose -f deployment/multi-node-small/docker/docker-compose.yml up -d
./add-team.sh example
dynamicbeat/dynamicbeat -e
```
