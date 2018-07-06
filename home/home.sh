#!/bin/bash
docker rmi -f smasher164/home:prod
docker rm -f home
docker run --name="home" -p 8082:8080 smasher164/home:prod