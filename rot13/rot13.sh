#!/bin/bash
docker rmi -f smasher164/rot13:prod
docker rm -f rot13
docker run --name="rot13" -p 8088:8080 smasher164/rot13:prod