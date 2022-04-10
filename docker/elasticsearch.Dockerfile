FROM docker.elastic.co/elasticsearch/elasticsearch:8.1.2

# Default password that should changed
ENV ELASTIC_PASSWORD=changeme

# Used to persist data if the container is stopped
VOLUME /usr/share/elasticsearch/data
# Used to share certificates between containers
VOLUME /usr/share/elasticsearch/config/certs

COPY configs/elasticsearch.yml config/elasticsearch.yml
COPY scripts/elasticsearch-entrypoint.sh /opt/entrypoint.sh
COPY scripts/elasticsearch-healthcheck.sh /opt/healthcheck.sh

EXPOSE 9200/tcp
HEALTHCHECK --interval=10s --timeout=10s --start-period=2m --retries=3 CMD [ "/opt/healthcheck.sh" ]
ENTRYPOINT [ "/bin/tini", "--" ]
CMD [ "/opt/entrypoint.sh" ]