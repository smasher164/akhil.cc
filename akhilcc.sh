#!/bin/bash
docker rmi -f smasher164/akhilcc:prod
docker rm -f akhilcc
docker run  --net="host" --name="akhilcc" -v /certs:/certs smasher164/akhilcc:prod