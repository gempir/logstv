


.PHONY: bot
bot:
	@echo "=== Building Bot ==="
	@docker build -t gempir/logstvbot bot

.PHONY: push_bot
push_bot: bot
	@echo "=== Pushing Bot to Dockerhub ==="
	@docker push gempir/logstvbot

.PHONY: deploy
deploy:
	@echo "=== Deploying compose files ==="
	@scp docker-compose.yml root@eros.logs.tv:/root
	@scp docker-compose.prod.yml root@eros.logs.tv:/root
	@ssh root@eros.logs.tv docker-compose -f docker-compose.yml -f docker-compose.prod.yml pull
	@ssh root@eros.logs.tv docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

.PHONY: release
release: bot push_bot deploy

provision: 
	ansible-playbook -i ansible/hosts ansible/playbook.yml --ask-vault-pass ${ARGS}

remove_all:
	docker rm `docker ps -aq`