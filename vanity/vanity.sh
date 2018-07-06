#!/bin/bash
docker rmi -f smasher164/vanity:prod
docker rm -f vanity
docker run --name="vanity" -p 8081:8080 smasher164/vanity:prod