###############################################################################
FROM golang:1.13.10 as ci
###############################################################################

RUN apt-get update

# Set up non-root user ########################################################

ARG USERNAME=scorestack
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Add non-root user
RUN groupadd --gid $USER_GID $USERNAME
RUN useradd -s /bin/bash --uid $USER_UID --gid $USER_GID -m $USERNAME

# Add sudo privileges to non-root user
RUN apt-get install -y sudo
RUN echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME
RUN chmod 0440 /etc/sudoers.d/$USERNAME

# Set up non-root user gopath
RUN mkdir -p /home/$USERNAME/go/src/github.com/s-newman
RUN chown -R $USER_UID:$USER_GID /home/$USERNAME/go

# Install Packages ############################################################

# Install build dependencies
RUN apt-get install -y \
    python-pip \
    virtualenv \
    git

# Clone go-elasticsearch repository
RUN git clone https://github.com/elastic/go-elasticsearch.git $GOPATH/src/github.com/elastic/go-elasticsearch
RUN cd $GOPATH/src/github.com/elastic/go-elasticsearch && git checkout v7.5.0

###############################################################################
FROM ci as devcontainer
###############################################################################

# Install packages ############################################################

# Install Go tools
RUN go get -v golang.org/x/tools/...

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.25.1

# Install a bunch of packages that vscode wants. I don't really know what all
# of these are, but they make the go extension work properly.
RUN GO111MODULE=on go get -v \
    honnef.co/go/tools/...@latest \
    golang.org/x/tools/cmd/gorename@latest \
    golang.org/x/tools/cmd/goimports@latest \
    golang.org/x/tools/cmd/guru@latest \
    golang.org/x/lint/golint@latest \
    github.com/mdempsky/gocode@latest \
    github.com/cweill/gotests/...@latest \
    github.com/haya14busa/goplay/cmd/goplay@latest \
    github.com/sqs/goreturns@latest \
    github.com/josharian/impl@latest \
    github.com/davidrjenni/reftools/cmd/fillstruct@latest \
    github.com/uudashr/gopkgs/...  \
    github.com/ramya-rao-a/go-outline@latest  \
    github.com/acroca/go-symbols@latest  \
    github.com/godoctor/godoctor@latest  \
    github.com/rogpeppe/godef@latest  \
    github.com/zmb3/gogetdoc@latest \
    github.com/fatih/gomodifytags@latest  \
    github.com/mgechev/revive@latest  \
    github.com/go-delve/delve/cmd/dlv@latest 2>&1
RUN go get -v github.com/alecthomas/gometalinter 2>&1
RUN go get -x -d github.com/stamblerre/gocode 2>&1
RUN go build -o $GOPATH/bin/gocode-gomod github.com/stamblerre/gocode