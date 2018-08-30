#!/usr/bin/env bash
set -e
set -u

BUILD_API=$1
BUILD_BOT=$2
BUILD_RELAYBROKER=$3

# Upload
scp ${BUILD_API} root@eros.logs.tv:/home/logstvapi
scp ${BUILD_BOT} root@eros.logs.tv:/home/logstvbot
scp ${BUILD_RELAYBROKER} root@eros.logs.tv:/home/logstvbot

# Extract
ssh root@eros.logs.tv tar xvf /home/logstvapi/build_api.tar.gz -C /home/logstvapi/
ssh root@eros.logs.tv tar xvf /home/logstvbot/build_bot.tar.gz -C /home/logstvbot/
ssh root@eros.logs.tv tar xvf /home/logstvbot/build_relaybroker.tar.gz -C /home/logstvbot/