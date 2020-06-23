###############################################################################
FROM node:10.19.0 as ci
# ARGS: KIBANA_VERSION, USER_GID, USER_UID, USERNAME
###############################################################################

RUN apt-get update

# Set up non-root user ########################################################
# The node container already comes with a "node" user of UID 1000, so we'll
# just use that.

ARG USERNAME=node
ARG USER_UID=1000
ARG USER_GID=${USER_UID}

# Add sudo privileges to non-root user
RUN apt-get install -y sudo
RUN echo ${USERNAME} ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/${USERNAME}
RUN chmod 0440 /etc/sudoers.d/${USERNAME}

# Install Packages ############################################################

# Install build dependencies
RUN apt-get install -y \
    git \
    libnss3
RUN npm install -g \
    eslint

# Clone correct version of Kibana
ARG KIBANA_VERSION=7.7.1
RUN git clone -b v${KIBANA_VERSION} --depth 1 https://github.com/elastic/kibana.git /home/${USERNAME}/kibana

# Set up plugin directory
RUN mkdir -p /home/${USERNAME}/kibana/plugins
RUN chown -R ${USER_UID}:${USER_GID} /home/${USERNAME}/kibana

# Bootstrap Kibana and install the node dependencies
RUN cd /home/${USERNAME}/kibana && sudo -u node yarn kbn bootstrap

###############################################################################
FROM ci as devcontainer
# ARGS: USERNAME
###############################################################################

ARG USERNAME=node

# Set up cluster
COPY files/kibana.dev.yml /home/${USERNAME}/kibana/config/kibana.dev.yml