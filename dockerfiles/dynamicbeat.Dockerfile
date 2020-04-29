FROM golang:1.13.10 as build

# Install build dependencies
RUN apt-get update && apt-get install -y \
    python-pip \
    virtualenv \
    git

# Clone go-elasticsearch repository
RUN git clone https://github.com/elastic/go-elasticsearch.git $GOPATH/src/github.com/elastic/go-elasticsearch
RUN cd $GOPATH/src/github.com/elastic/go-elasticsearch && git checkout v7.5.0

FROM build as devcontainer

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID