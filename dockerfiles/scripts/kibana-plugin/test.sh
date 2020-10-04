#!/bin/bash
set -euxo pipefail

cd $HOME/kibana/plugins/scorestack
yarn plugin-helpers build