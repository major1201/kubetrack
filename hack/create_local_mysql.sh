#!/usr/bin/env bash

set -euo pipefail

docker run -d --name ktpg -p 5432:5432 \
    -e POSTGRES_PASSWORD=password \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v "${HOME}/Downloads/pgdata":/var/lib/postgresql/data \
    postgres
