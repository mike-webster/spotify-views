
.PHONY: serve_prod
serve_prod:
	nohup ./spotify-views > spotify-views.log 2>&1 &

.PHONY: kill_prod
kill_prod:
	kill -9 $(lsof -i tcp:3000 | tail -n +2 | awk '{print $2}')

.PHONY: find_pid
find_pid:
	pgrep -a spotify-views

.PHONY: build
build: 
	go build -o spotify-views cmd/spotify-views/main.go

.PHONY: reset
reset: build kill_prod serve_prod