FROM ubuntu:bionic

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Add a non-root user to the container
ARG USERNAME=scorestack
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Configure apt and install some common tools
RUN apt-get update && \
    apt-get -y install --no-install-recommends apt-utils dialog 2>&1 && \
    apt-get -y install git openssh-client iproute2 procps lsb-release vim virtualenv

# Create non-root user
RUN groupadd --gid $USER_GID $USERNAME && \
    useradd -s /bin/bash --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    apt-get install -y sudo && \
    echo $USERNAME ALL\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

# Install golang
RUN add-apt-repository -y ppa:longsleep/golang-backports && \
    apt-get update && \
    apt-get install -y golang-go

# Install golang tools
RUN mkdir -p /tmp/gotools && \
    cd /tmp/gotools && \
    GO111MODULE=on go get -v golang.org/x/tools/gopls@latest 2>&1 && \
    GO111MODULE=on go get -v && \
        honnef.co/go/tools/...@latest && \
        golang.org/x/tools/cmd/gorename@latest && \
        golang.org/x/tools/cmd/goimports@latest && \
        golang.org/x/tools/cmd/guru@latest && \
        golang.org/x/lint/golint@latest && \
        github.com/mdempsky/gocode@latest && \
        github.com/cweill/gotests/...@latest && \
        github.com/haya14busa/goplay/cmd/goplay@latest && \
        github.com/sqs/goreturns@latest && \
        github.com/josharian/impl@latest && \
        github.com/davidrjenni/reftools/cmd/fillstruct@latest && \
        github.com/uudashr/gopkgs/cmd/gopkgs@latest  && \
        github.com/ramya-rao-a/go-outline@latest  && \
        github.com/acroca/go-symbols@latest  && \
        github.com/godoctor/godoctor@latest  && \
        github.com/rogpeppe/godef@latest  && \
        github.com/zmb3/gogetdoc@latest && \
        github.com/fatih/gomodifytags@latest  && \
        github.com/mgechev/revive@latest  && \
        github.com/go-delve/delve/cmd/dlv@latest 2>&1 \
    go get -v github.com/alecthomas/gometalinter 2>&1 && \
    go get -x -d github.com/stamblerre/gocode 2>&1 && \
    go build -o gocode-gomod github.com/stamblerre/gocode && \
    mv gocode-gomod $GOPATH/bin/ && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin 2>&1

# Set up non-root user gopath
RUN mkdir -p /home/$USERNAME/go/src/github.com/s-newman && \
    chown -R $USER_UID:$USER_GID /home/$USERNAME/go

# Clean up
RUN apt-get autoremove -y && \
    apt-get clean -y && \
    rm -rf /var/lib/apt/lists/* /go/src /tmp/gotools

# Update this to "on" or "off" as appropriate
ENV GO111MODULE=auto

# Switch back to dialog for any ad-hoc use of apt-get
ENV DEBIAN_FRONTEND=dialog

# Set go path envs in here just in case they don't get set externally
ENV GOPATH=/home/vscode/go
ENV PATH=${PATH}:${GOPATH}/bin