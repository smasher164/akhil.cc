#!/bin/bash
docker rmi -f smasher164/staging:prod
docker rm -f staging
docker run --name="staging" -p 8086:8080 smasher164/staging:prod