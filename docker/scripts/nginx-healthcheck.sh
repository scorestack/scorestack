#!/bin/bash
set -eou pipefail

curl --slient http://localhost:8000/health | grep -q "healthy"