default: bot api 

push: push_bot push_api

.PHONY: bot
bot:
	@echo "=== Building Bot ==="
	@docker build -t gempir/logstvbot bot

.PHONY: api
api:
	@echo "=== Building Api ==="
	@docker build -t gempir/logstvapi api

.PHONY: push_bot
push_bot: bot
	@echo "=== Pushing Bot to Dockerhub ==="
	@docker push gempir/logstvbot

.PHONY: push_api
push_api: api
	@echo "=== Pushing Api to Dockerhub ==="
	@docker push gempir/logstvapi

.PHONY: deploy
deploy:
	@echo "=== Deploying compose files ==="
	@scp docker-compose.yml root@apa.logs.tv:/root
	@scp docker-compose.prod.yml root@apa.logs.tv:/root
	@ssh root@apa.logs.tv docker-compose -f docker-compose.yml -f docker-compose.prod.yml pull
	@ssh root@apa.logs.tv docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

.PHONY: release
release: bot push_bot deploy

provision: 
	ansible-playbook -i ansible/hosts ansible/playbook.yml --ask-vault-pass ${ARGS}

remove_all:
	docker rm `docker ps -aq`