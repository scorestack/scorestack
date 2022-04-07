#!/bin/bash
set -eou pipefail

curl -s --cacert -I https://localhost:5601 | grep -q 'HTTP/1.1 302 Found'