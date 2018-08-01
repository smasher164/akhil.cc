#!/bin/bash
docker rmi -f smasher164/static:prod
docker rm -f static
docker run --name="static" -p 8083:8080 smasher164/static:prod