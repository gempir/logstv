#!/usr/bin/env bash
set -e
set -u

BUILD_API=$1
BUILD_BOT=$2
BUILD_RELAYBROKER=$3

echo "=== Upload services ==="
scp ${BUILD_API} root@eros.logs.tv:/home/logstvapi
scp ${BUILD_BOT} root@eros.logs.tv:/home/logstvbot
scp ${BUILD_RELAYBROKER} root@eros.logs.tv:/home/logstvbot

echo "=== Extracting services ==="
ssh root@eros.logs.tv tar xvf /home/logstvapi/build_api.tar.gz -C /home/logstvapi/
ssh root@eros.logs.tv tar xvf /home/logstvbot/build_bot.tar.gz -C /home/logstvbot/
ssh root@eros.logs.tv tar xvf /home/logstvbot/build_relaybroker.tar.gz -C /home/logstvbot/

echo "=== Restarting services ==="
ssh root@eros.logs.tv service relaybroker restart
ssh root@eros.logs.tv service logstvbot restart
ssh root@eros.logs.tv service logstvapi restart