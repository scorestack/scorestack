#!/bin/bash
set -eou pipefail

curl --silent http://localhost:8000/health | grep -q "healthy"