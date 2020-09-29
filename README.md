# scorestack

![Dynamicbeat](https://github.com/scorestack/scorestack/workflows/Dynamicbeat/badge.svg)

A security competition scoring system built on the Elastic stack.

## Building dynamicbeat

Run the following command to build and test dynamicbeat with docker:

```shell
docker-compose -f dockerfiles/docker-compose.yml run dynamicbeat-ci /scripts/dynamicbeat-test.sh
```

The compiled binary can be found at `dynamicbeat/dynamicbeat`.

## Running everything

After dynamicbeat has been built, run the following commands to set up an
instance of Scorestack with the example check configurations loaded up and run
an instance of dynamicbeat against those checks.

```shell
sudo sysctl -w vm.max_map_count=262144
docker-compose -f deployment/multi-node-small/docker/docker-compose.yml up -d
./add-team.sh example
dynamicbeat/dynamicbeat -e
```
