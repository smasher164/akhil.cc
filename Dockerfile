FROM octoblu/alpine-ca-certificates
EXPOSE 80
EXPOSE 443
VOLUME ["/certs"]
COPY akhilcc /
COPY conf.toml /
ENTRYPOINT ["/akhilcc", "--conf", "/conf.toml", "--site", "www.akhil.cc"]
