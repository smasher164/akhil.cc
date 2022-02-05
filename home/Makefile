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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o home
	for file in $^ ; do \
		echo $${file::$${#file}-3} | xargs -I {} mexdown html {}.xd -o {}.html ; \
	done
	docker build . --tag home:prod
	docker tag home:prod smasher164/home:prod
	docker push smasher164/home:prod
	ssh core@akhil.cc "sudo systemctl restart home"