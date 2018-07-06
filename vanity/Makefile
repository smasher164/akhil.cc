build:
	docker ps -a | awk '{ print $$1,$$2 }' | grep vanity:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	#docker rmi -f vanity:prod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build . --tag vanity:prod
	docker run --rm -detach -p 8080:8080 vanity:prod
deploy:
	docker ps -a | awk '{ print $$1,$$2 }' | grep vanity:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build . --tag vanity:prod | tail -n 2 | head -n 1 | awk '{print $$3}' | xargs -I {} docker tag {} smasher164/vanity:prod
	docker push smasher164/vanity:prod
	ssh core@akhil.cc "docker pull smasher164/vanity:prod && sudo systemctl restart vanity"