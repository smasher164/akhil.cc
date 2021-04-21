FROM ubuntu:16.04
EXPOSE 8080
RUN apt-get update && apt-get install -y wget perl && rm -rf /var/lib/apt/lists/*
COPY blog $HOME
COPY post $HOME
COPY svglatex $HOME
COPY posts $HOME/posts
COPY static $HOME/static
COPY posts.toml $HOME
ENV PATH="/:$PATH"
ENTRYPOINT ["blog", "--conf", "posts.toml", "--posts", "posts", "--static", "static"]