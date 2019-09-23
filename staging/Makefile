deploy:
	docker ps -a | awk '{ print $$1,$$2 }' | grep staging:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build . --tag staging:prod | tail -n 2 | head -n 1 | awk '{print $$3}' | xargs -I {} docker tag {} smasher164/staging:prod
	docker push smasher164/staging:prod
	ssh core@akhil.cc "docker pull smasher164/staging:prod && sudo systemctl restart staging"