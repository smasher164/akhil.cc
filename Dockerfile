FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates
EXPOSE 80
EXPOSE 443
VOLUME ["/certs"]
COPY akhilcc /
COPY conf.toml /
ENTRYPOINT ["/akhilcc", "--conf", "/conf.toml", "--site", "www.akhil.cc"]
