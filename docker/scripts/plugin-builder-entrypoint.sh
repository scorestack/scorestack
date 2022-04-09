#!/bin/bash
set -eou pipefail

if [ -f plugins/scorestack/build/scorestack-8.1.2.zip ]; then
  echo "Plugin already built, exiting"
  exit 0
fi

if [ ! -d .git ]; then
  echo "Cloning Kibana repo"
  git init
  git remote add origin https://github.com/elastic/kibana.git
  git fetch --depth 1 origin tag v8.1.2
  git checkout v8.1.2
fi

echo "Performing initial bootstrapping"
cd plugins/scorestack
yarn kbn bootstrap

echo "Building plugin"
yarn build

echo "Plugin built!"