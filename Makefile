

build: build_api build_bot build_relaybroker

.PHONY: build_api
build_api:
	@echo "Building API"
	@cd api && go get ./... && env GOOS=linux GOARCH=amd64 go build
	@tar -czvf build_api.tar.gz api/api

.PHONY: build_bot
build_bot:
	@echo "Building Bot"
	@cd bot && go get ./... && env GOOS=linux GOARCH=amd64 go build
	@tar -czvf build_bot.tar.gz bot/bot

.PHONY: build_relaybroker
build_relaybroker: clone_relaybroker
	@echo "Building Relaybroker"
	@cd relaybroker && go get ./... && env GOOS=linux GOARCH=amd64 go build
	@tar -czvf build_relaybroker.tar.gz relaybroker/relaybroker

.PHONY: deploy
deploy: 
	@echo "Deploying 3 Apps"
	@deploy/deploy.sh build_api.tar.gz build_bot.tar.gz	build_relaybroker.tar.gz

provision: 
	ansible-playbook -i ansible/hosts ansible/playbook.yml --ask-vault-pass

clone_relaybroker:
	@echo "Cloning Relaybroker"
	@rm -rf relaybroker
	@git clone https://github.com/gempir/relaybroker relaybroker