HOST := localhost

GO_ENV ?= development

APP_NAME := spotify-views
USER := webby

NET_NAME := sv-net

CLIENT_NAME := sv-client
CLIENT_PORT := 3000

DB_NAME := sv-db
DB_PORT := 3306

RED_NAME := redis
REDIS_PORT := 6379

API_NAME := sv-api
API_PORT := 3001

RELEASE_DATE ?= 
RELEASE_DIR := /home/spotify-views/
API_RELEASE_DIR := releases-api

## Encryption Util
.PHONY: build_enc
build_enc:
	go build -o enc ./cmd/encrypt 

.PHONY: enc
enc: build_enc
	GO_ENV=$(GO_ENV) ./enc

.PHONY: refresh_secrets
refresh_secrets: build_enc
	./enc -e -in=secrets.yaml -out=secrets.enc

## Local dev tools
.PHONY: dev
dev:
	go build -o $(APP_NAME) ./cmd/spotify-views/main.go 
	HOST=$(HOST) PORT=$(API_PORT) GO_ENV=$(GO_ENV) ./$(APP_NAME)

.PHONY: check
check:
	go build -o $(APP_NAME) ./cmd/spotify-views/main.go
	HOST=$(HOST) PORT=$(API_PORT) GO_ENV=$(GO_ENV) ./$(APP_NAME) -check

## Production tools
.PHONY: serve_prod
serve_prod:
	nohup PORT=$(API_PORT) ./$(APP_NAME) > $(APP_NAME).log 2>&1 &

.PHONY: kill_prod
kill_prod:
	kill -9 $(lsof -i tcp:$(API_PORT) | tail -n +2 | awk '{print $2}')

.PHONY: find_pid
find_pid:
	pgrep -a $(APP_NAME)

.PHONY: release_api
release_api: test refresh_secrets
	echo $(RELEASE_DATE)
ifeq ($(RELEASE_DATE),)
	echo "cannot release without a date provided"
	exit 1
endif

ifeq ($(GO_ENV),production)
	GO_ENV=$(GO_ENV) GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) ./cmd/spotify-views/main.go 

	# create the release directory
	ssh -i ~/.ssh/id_rsa root@spotify-views.com mkdir $(RELEASE_DIR)$(API_RELEASE_DIR)/$(RELEASE_DATE)

	# copy the binary into the new release directory
	scp -r ./$(APP_NAME) root@spotify-views.com:$(RELEASE_DIR)$(API_RELEASE_DIR)/$(RELEASE_DATE)/$(APP_NAME)

	# copy secrets with the new release
	scp -r ./secrets.enc root@spotify-views.com:$(RELEASE_DIR)$(API_RELEASE_DIR)/$(RELEASE_DATE)/secrets.enc

	# copy the new release into the live directory, overwriting what's already there
	ssh root@spotify-views.com cp -rf $(RELEASE_DIR)$(API_RELEASE_DIR)/$(RELEASE_DATE)/$(APP_NAME) $(RELEASE_DIR)$(API_RELEASE_DIR)/live
else
	echo "cannot release unless GO_ENV set to 'production'"
endif


## Docker tools

## ## ## start - this will clear any existing containers, build the app,
## ## ##         start the db container, and start the react app.
.PHONY: start
start: clear build init_db run client_start

.PHONY: client_logs
client_logs:
	docker logs $(CLIENT_NAME)

.PHONY: api_logs
api_logs:
	docker logs $(API_NAME) -n 1000

.PHONY: clear_app
clear_app:
	docker container rm $(API_NAME) -f

.PHONY: build
build: 
	docker build . -t $(APP_NAME) \
		--build-arg HOST=$(HOST) \
		--build-arg PORT=$(API_PORT) \
		--build-arg MASTER_KEY="$(MASTER_KEY)" \
		--build-arg GO_ENV=$(GO_ENV) 
	
.PHONY: run
run: clear_app build
	docker run \
		-d \
		-p $(API_PORT):$(API_PORT) \
		--name $(API_NAME) \
		--network $(NET_NAME) \
		$(APP_NAME) 

.PHONY: in_api
in_api:
	docker exec -it $(API_NAME) sh

.PHONY: clear_db
clear_db:
	docker volume prune -f
	docker container rm $(DB_NAME) -f

.PHONY: clear_network
clear_network:
	docker network rm $(NET_NAME)

.PHONY: create_network
create_network: 
	docker network create $(NET_NAME)

.PHONY: clear
clear: clear_app clear_db client_clear clear_network 

.PHONY: start_db
start_db: 
	docker container rm $(DB_NAME) -f
	docker pull mysql
	docker run \
		-p $(DB_PORT):$(DB_PORT) \
		--name $(DB_NAME) \
		--volume=/Users/$(USER)/storage/docker/mysql-data:/var/lib/mysql \
		--network $(NET_NAME) \
		-e MYSQL_ROOT_PASSWORD=pass \
		-d \
		mysql
	@sleep 15s
	
.PHONY: init_db
init_db: clear_db start_db
	@docker exec -i $(DB_NAME) \
		mysql -uroot -ppass --protocol=tcp -h localhost -P $(DB_PORT) < ./data/create_db.sql

.PHONY: start_redis
start_redis:
	docker container rm $(RED_NAME) -f
	docker pull redis
	docker run \
		--name $(RED_NAME) \
		--volume /Users/$(USER)/storage/redis/redis-data:/data \
		--network $(NET_NAME) \
		-p $(REDIS_PORT):$(REDIS_PORT) \
		-d \
		redis

## Testing

.PHONY: test
test: 
	GO_ENV=test go test ./... -cover

.PHONY: convey
convey:
	sudo goconvey -excludedDirs="web,cmd,client" -port=8888

## Release
.PHONY: release
release: 
	git pull 
	go build ./cmd/spotify-views
	$(MAKE) kill_prod
	$(MAKE) serve_prod

## Client
.PHONY: client_build
client_build:
	docker build -f ./client/Dockerfile ./client/ -t $(CLIENT_NAME)

.PHONY: client_clear
client_clear:
	docker container rm $(CLIENT_NAME) -f

.PHONY: client_start
client_start: client_build client_clear
	docker run \
		--network $(NET_NAME) \
		--name $(CLIENT_NAME) \
		-d \
		-p $(CLIENT_PORT):$(CLIENT_PORT) \
		-v /Users/webby/Code/spotify-views/client:/app \
		$(CLIENT_NAME)