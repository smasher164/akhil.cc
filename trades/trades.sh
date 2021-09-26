#!/bin/bash
docker rmi -f smasher164/trades:prod
docker rm -f trades
docker run --name="trades" -p 8087:8080 smasher164/trades:prod