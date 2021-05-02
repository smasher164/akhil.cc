run:
	docker ps -a | awk '{ print $$1,$$2 }' | grep akhilcc:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	docker rmi -f akhilcc:prod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o akhilcc
	docker build . --tag akhilcc:prod
	docker run --rm --detach akhilcc:prod
build:
	docker ps -a | awk '{ print $$1,$$2 }' | grep akhilcc:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	docker rmi -f akhilcc:prod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o akhilcc
	docker build . --tag akhilcc:prod
deploy:
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o akhilcc
	docker build . --tag akhilcc:prod
	docker tag akhilcc:prod smasher164/akhilcc:prod
	docker push smasher164/akhilcc:prod
	ssh core@akhil.cc "sudo systemctl restart akhilcc"
mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o akhilcc
