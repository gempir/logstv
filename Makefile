

build: build_api build_bot

.PHONY: build_api
build_api:
	@echo "Building API"
	@cd api && go get ./... && env GOOS=linux GOARCH=amd64 go build
	@tar -czvf api_build.tar.gz api/api

.PHONY: build_bot
build_bot:
	@echo "Building Bot"
	@cd bot && go get ./... && env GOOS=linux GOARCH=amd64 go build
	@tar -czvf bot_build.tar.gz bot/bot

.PHONY: deploy
deploy: 
	@echo "Deploying 2 Apps"
	@deploy/deploy.sh api_build.tar.gz bot_build.tar.gz	

provision: 
	ansible-playbook -i ansible/hosts ansible/playbook.yml --ask-vault-pass