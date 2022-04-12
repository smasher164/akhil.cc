build:
	docker ps -a | awk '{ print $$1,$$2 }' | grep static:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	#docker rmi -f static:prod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build . --tag static:prod
	docker run --rm -detach -p 8080:8080 static:prod
deploy:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build . --tag static:prod
	docker tag static:prod smasher164/static:prod
	docker push smasher164/static:prod
	ssh core@akhil.cc "docker pull smasher164/static:prod && sudo systemctl restart static"