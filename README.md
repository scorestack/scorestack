# scorestack

A security competition scoring system built on the Elastic stack.

## Building dynamicbeat

For a number of reasons, we use vscode dev container to build dynamicbeat. Once
you open the project in the provided dev container, run the following commands
to build the beat:

```shell
sudo chown vscode:vscode -R ~/go
cd dynamicbeat
make setup
go get ./...
mage build
```