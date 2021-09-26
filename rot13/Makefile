build:
	docker ps -a | awk '{ print $$1,$$2 }' | grep rot13:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	#docker rmi -f rot13:prod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build . --tag rot13:prod
	docker run --rm -detach -p 8080:8080 rot13:prod
deploy:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build . --tag rot13:prod
	docker tag rot13:prod smasher164/rot13:prod
	docker push smasher164/rot13:prod
	ssh core@akhil.cc "sudo systemctl restart rot13"