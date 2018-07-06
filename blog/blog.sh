#!/bin/bash
docker rmi -f smasher164/blog:prod
docker rm -f blog
docker run --name="blog" -p 8080:8080 smasher164/blog:prod