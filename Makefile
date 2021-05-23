HOST := localhost
PORT := 8081
GO_ENV := development
APP_NAME := spotify-views

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

.PHONY: clear
clear:
	docker container rm sv-dev -f

.PHONY: build
build: 
	docker build . -t $(APP_NAME) --no-cache \
	--build-arg HOST=$(HOST) \
	--build-arg PORT=$(PORT)
	
.PHONY: run
run: clear build
	docker run \
	-p 8080:8080 \
	--name sv-dev \
	-v ~/mike-webster/spotify-views:/app \
	spotify-views 

## Testing

.PHONY: test
test: 
	GO_ENV=test go test ./... -cover

.PHONY: convey
convey:
	sudo goconvey -excludedDirs="web,cmd,client" -port=8888
