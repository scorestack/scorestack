#!/bin/bash
set -euxo pipefail

cd $HOME/kibana/plugins/scorestack
yarn kbn bootstrap
yarn prebuild
yarn build