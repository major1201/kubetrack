#!/usr/bin/env bash

set -euo pipefail

docker run -d --name ktmysql -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=password \
    -v "${HOME}/Downloads/mysqldata":/var/lib/mysql \
    mysql
