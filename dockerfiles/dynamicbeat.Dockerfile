###############################################################################
FROM golang:1.20 as ci
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
RUN mkdir -p /home/$USERNAME/go
RUN chown -R $USER_UID:$USER_GID /home/$USERNAME/go

# Set up target directory
RUN mkdir -p /home/$USERNAME/scorestack
RUN chown -R $USER_UID:$USER_GID /home/$USERNAME/scorestack

# Install Packages ############################################################

# Install build dependencies
RUN apt-get install -y \
    python3-pip \
    virtualenv \
    git

# Install docker CLI
RUN apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
RUN apt-get update && apt-get install -y docker-ce-cli

###############################################################################
FROM ci as devcontainer
###############################################################################

# Install packages ############################################################

# Install Go tools
RUN go install -v golang.org/x/tools/...@latest 2>&1

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2

# Install a bunch of packages that vscode wants. I don't really know what all
# of these are, but they make the go extension work properly.
RUN go install -v github.com/cweill/gotests/gotests@v1.6.0 2>&1
RUN go install -v github.com/fatih/gomodifytags@v1.16.0 2>&1
RUN go install -v github.com/josharian/impl@v1.1.0 2>&1
RUN go install -v github.com/haya14busa/goplay/cmd/goplay@v1.0.0 2>&1
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest 2>&1
RUN go install -v honnef.co/go/tools/cmd/staticcheck@latest 2>&1
RUN go install -v golang.org/x/tools/gopls@latest 2>&1
