#!/bin/bash
set -eou pipefail

if [ ! -d .git ]; then
  echo "Cloning Kibana repo"
  git init
  git remote add origin https://github.com/elastic/kibana.git
  git fetch --depth 1 origin tag v8.1.2
  git checkout v8.1.2

  echo "Performing initial bootstrapping"
  yarn kbn bootstrap
fi

yarn start --dev --allow-root