HOST := localhost
PORT := 8081
GO_ENV := development
APP_NAME := spotify-views
USER := webby

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
	HOST=$(HOST) PORT=$(PORT) GO_ENV=$(GO_ENV) ./$(APP_NAME)

## Production tools
.PHONY: serve_prod
serve_prod:
	nohup ./$(APP_NAME) > $(APP_NAME).log 2>&1 &

.PHONY: kill_prod
kill_prod:
	kill -9 $(lsof -i tcp:3000 | tail -n +2 | awk '{print $2}')

.PHONY: find_pid
find_pid:
	pgrep -a $(APP_NAME)

## Docker tools

.PHONY: clear_app
clear:
	docker container rm sv-dev -f

.PHONY: build
build: 
	docker build . -t $(APP_NAME) \
		--build-arg HOST=$(HOST) \
		--build-arg PORT=$(PORT) \
		--build-arg MASTER_KEY="$(MASTER_KEY)" \
		--build-arg GO_ENV=development 
	
.PHONY: run
run: clear_app build
	docker run \
		-it \
		-p 8081:8081 \
		--name sv-dev \
		--network sv-net \
		$(APP_NAME) 

.PHONY: in
in:
	docker exec -it sv-dev sh

.PHONY: clear_db
clear_db:
	docker volume prune -f

.PHONY: clear_network
clear_network:
	docker network rm sv-net

.PHONY: clear
clear: clear_db clear_network clear_app

.PHONY: start_db
start_db:
	docker network create sv-net
	docker container rm sv-db -f
	docker pull mysql
	docker run \
		-p 3306:3306 \
		--name sv-db \
		--volume=/Users/$(USER)/storage/docker/mysql-data:/var/lib/mysql \
		--network sv-net \
		-e MYSQL_ROOT_PASSWORD=pass \
		-d \
		mysql
	@sleep 15s
	
.PHONY: init_db
init_db: clear_db start_db
	@docker exec -i sv-db \
		mysql -uroot -ppass --protocol=tcp -h localhost -P 3306 < ./data/create_db.sql

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