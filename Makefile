SHELL := /bin/bash

.PHONY: deploy-web
# Build admin for linux amd64 binary
deploy-web:
	@cd web && bash deployments/build.sh && expect deployments/push.sh $(USER) $(PWD) $(IP)

.PHONY: deploy-admin
# Deploy binary to remote linux server, e.g. make deploy-binary USER=root PWD=123456 IP=192.168.1.10
deploy-admin:
	@cd server && make deploy-binary USER=$(USER) PWD=$(PWD) IP=$(IP)


.PHONY: deploy
# Deploy binary to remote linux server, e.g. make deploy-binary USER=root PWD=123456 IP=192.168.1.10
deploy:
	@cd server && make deploy-binary USER=$(USER) PWD=$(PWD) IP=$(IP)
	@cd web && bash deployments/build.sh && expect deployments/push.sh $(USER) $(PWD) $(IP)

.PHONY: clean
# Clean binary file, cover.out, template file
clean:
	#rm -rf admin-binary.tar.gz
	@cd server && make clean

# Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[1;36m  %-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all
