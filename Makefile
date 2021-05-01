.PHONY: start
start:
	CLIENT_ID=TEST CLIENT_SECRET=TEST HOST='localhost:8080' LYRICS_KEY=bV9WV_Cs-AICvaY0uiCEuf4uaH74aJyHTAvX3Zr_BJA2CGgOm1njlrxWfYT_cYv- go run main.go handlers.go helpers.go
.PHONY: tail_prod
tail_prod: 
	heroku logs --tail -a spotify-views

.PHONY: serve_prod
serve_prod:
	nohup ./spotify-views > spotify-views.log 2>&1 &

.PHONY: find_pid
find_pid:
	pgrep -a spotify-views
