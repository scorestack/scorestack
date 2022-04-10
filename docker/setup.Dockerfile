FROM docker.elastic.co/elasticsearch/elasticsearch:8.1.2
USER root

# Used to share certificates between containers
VOLUME /usr/share/elasticsearch/config/certs

# Default passwords that should changed
ENV ELASTIC_PASSWORD=changeme
ENV KIBANA_PASSWORD=changeme

COPY configs/instances.yml config/certs/instances.yml
COPY scripts/setup-entrypoint.sh /opt/entrypoint.sh

ENTRYPOINT [ "/bin/tini", "--" ]
CMD [ "/opt/entrypoint.sh" ]