FROM node:16.14.2 AS plugin-builder
WORKDIR /usr/share/kibana

# Note this container must be built from the root of the repository so the
# kibana-plugin/ folder can be included in the build context. This differs from
# the rest of the images, which are built from within the docker/ folder to
# reduce the size of the build context.
COPY kibana-plugin /usr/share/kibana/plugins/scorestack
COPY docker/scripts/plugin-builder-entrypoint.sh /opt/entrypoint.sh

RUN /opt/entrypoint.sh

FROM docker.elastic.co/kibana/kibana:8.1.2

# Default password that should be changed
ENV ELASTICSEARCH_PASSWORD=changeme

# Used to persist data if the container is stopped
VOLUME /usr/share/kibana/data
# Used to share certificates between containers
VOLUME /usr/share/kibana/config/certs

COPY --from=plugin-builder /usr/share/kibana/plugins/scorestack/build/scorestack-8.1.2.zip /opt/plugin/build/scorestack-8.1.2.zip
COPY docker/configs/kibana.yml config/kibana.yml
COPY docker/scripts/kibana-entrypoint.sh /opt/entrypoint.sh
COPY docker/scripts/kibana-healthcheck.sh /opt/healthcheck.sh

EXPOSE 5601/tcp
HEALTHCHECK --interval=10s --timeout=10s --start-period=2m --retries=3 CMD [ "/opt/healthcheck.sh" ]
ENTRYPOINT [ "/bin/tini", "--" ]
CMD [ "/opt/entrypoint.sh" ]