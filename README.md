# scorestack

![Dynamicbeat](https://github.com/s-newman/scorestack/workflows/Dynamicbeat/badge.svg)

A security competition scoring system built on the Elastic stack.

## Running Code Generation

To manage type definitions for checks across multiple languages, we use OpenAPI (formerly known as Swagger API) code generation. All API definitions can be found under the `api/` directory.

Before writing any code or building any binaries, you must generate the code from the API definitions. Docker containers have been configured to simplify this process.

To generate the golang code for Dynamicbeat, run:

```shell
docker-compose -f dockerfiles/docker-compose.yml up api-generator-go
```

To generate the typescript code for the Kibana plugin, run:

```shell
docker-compose -f dockerfiles/docker-compose.yml up api-generator-typescript
```

## Building dynamicbeat

Run the following command to build and test dynamicbeat with docker:

```shell
docker-compose -f dockerfiles/docker-compose.yml run dynamicbeat-ci /scripts/dynamicbeat-test.sh
```

The compiled binary can be found at `dynamicbeat/dynamicbeat`.

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
