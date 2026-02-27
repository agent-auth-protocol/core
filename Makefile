start: # start the server locally
	go run main.go

docker-build: # build the optimized multi-stage docker image
	docker build -t agentauth-core:local .

docker-run: # run the server in a background container
	docker run -d -p 8080:8080 --name core-server agentauth-core:local

docker-stop: # stop and remove the running container
	docker rm -f core-server

docker-logs: # tail the logs of the running container
	docker logs -f core-server