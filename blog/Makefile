build:
	docker ps -a | awk '{ print $$1,$$2 }' | grep blog_1:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	#docker rmi -f blog:prod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/smasher164/post
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/smasher164/svglatex
	docker build . --tag blog_1:prod
	docker run --rm -p 8080:8080 blog_1:prod
deploy: posts/*.xd
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/smasher164/post
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/smasher164/svglatex
	for file in $^ ; do \
		echo $${file::$${#file}-3} | xargs -I {} mexdown html {}.xd -o {}.html ; \
	done
	docker build . --tag blog:prod
	docker tag blog:prod smasher164/blog:prod
	docker push smasher164/blog:prod
	ssh core@akhil.cc "sudo systemctl restart blog"