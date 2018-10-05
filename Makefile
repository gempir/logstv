

# build: build_api build_bot build_relaybroker

# .PHONY: build_api
# build_api:
# 	@echo "=== Building API ==="
# 	@cd api && go get ./... && env GOOS=linux GOARCH=amd64 go build
# 	@tar -czvf build_api.tar.gz api/api

.PHONY: build_bot
build_bot:
	@echo "=== Building Bot ==="
	@docker build -t gempir/logstvbot bot

.PHONY: dev_bot
dev_bot: build_bot
	@echo "=== Running Bot ==="
	@docker run --rm -p 8025:8025 -v `pwd`/bot/logs:/var/twitch_logs -v `pwd`/bot/channels:/etc/channels --name logstvbot gempir/logstvbot 

.PHONY: build_relaybroker
build_relaybroker: clone_relaybroker
	@echo "=== Building Relaybroker ==="
	@cd relaybroker && go get ./... && env GOOS=linux GOARCH=amd64 go build
	@tar -czvf build_relaybroker.tar.gz relaybroker/relaybroker

.PHONY: deploy
deploy: 
	@echo "=== Deploying 3 Apps ==="
	@deploy/deploy.sh build_api.tar.gz build_bot.tar.gz	build_relaybroker.tar.gz

provision: 
	ansible-playbook -i ansible/hosts ansible/playbook.yml --ask-vault-pass ${ARGS}

remove_all:
	docker rm `docker ps -aq`