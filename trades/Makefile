deploy:
	docker ps -a | awk '{ print $$1,$$2 }' | grep trades:prod | awk '{print $$1 }' | xargs -I {} docker rm -f {}
	cd ~/smasher164/trades && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build && docker build . --tag trades:prod | tail -n 2 | head -n 1 | awk '{print $$3}' | xargs -I {} docker tag {} smasher164/trades:prod
	docker push smasher164/trades:prod
	ssh core@akhil.cc "docker pull smasher164/trades:prod && sudo systemctl restart trades"