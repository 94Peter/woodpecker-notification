local-build:
	docker build -t mobydeck/ci-teams-notification .
	docker image prune -f


run: 
	go run main.go
