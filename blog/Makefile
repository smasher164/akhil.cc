build:
	docker ps -a | awk '{ print $$1,$$2 }' | grep blog_1:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	#docker rmi -f blog:prod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/smasher164/post
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/smasher164/svglatex
	docker build . --tag blog_1:prod
	docker run --rm -detach -p 8080:8080 blog_1:prod
deploy:
	docker ps -a | awk '{ print $$1,$$2 }' | grep blog_1:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/smasher164/post
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/smasher164/svglatex
	docker build . --tag blog_1:prod | tail -n 2 | head -n 1 | awk '{print $$3}' | xargs -I {} docker tag {} smasher164/blog:prod
	docker push smasher164/blog:prod
	ssh core@akhil.cc "docker pull smasher164/blog:prod && sudo systemctl restart blog"