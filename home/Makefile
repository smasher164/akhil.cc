build: *.xd
	docker ps -a | awk '{ print $$1,$$2 }' | grep home:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	#docker rmi -f home:prod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	for file in $^ ; do \
		echo $${file::$${#file}-3} | xargs -I {} mexdown html {}.xd -o {}.html ; \
	done
	docker build . --tag home:prod
	docker run -detach -p 8080:8080 home:prod
deploy: *.xd
	docker ps -a | awk '{ print $$1,$$2 }' | grep home:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	for file in $^ ; do \
		echo $${file::$${#file}-3} | xargs -I {} mexdown html {}.xd -o {}.html ; \
	done
	docker build . --tag home:prod | tail -n 2 | head -n 1 | awk '{print $$3}' | xargs -I {} docker tag {} smasher164/home:prod
	docker push smasher164/home:prod
	ssh core@akhil.cc "docker pull smasher164/home:prod && sudo systemctl restart home"