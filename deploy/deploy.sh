#!/usr/bin/env bash
set -e
set -u

BUILD_API=$1
BUILD_BOT=$2


# Upload
scp ${BUILD_API} root@eros.logs.tv:/home/logstvapi
scp ${BUILD_BOT} root@eros.logs.tv:/home/logstvbot

# Extract
ssh root@eros.logs.tv tar xvf /home/logstvapi/api_build.tar.gz -C /home/logstvapi/
ssh root@eros.logs.tv tar xvf /home/logstvbot/bot_build.tar.gz -C /home/logstvbot/