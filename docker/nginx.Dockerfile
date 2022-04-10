FROM nginx

# Used to share certificates between containers. Add your certificate and key
# here if you want to use a custom HTTPS certificate.
VOLUME /etc/nginx/certs

COPY docker/configs/proxy.conf /etc/nginx/conf.d/default.conf
COPY docker/scripts/nginx-entrypoint.sh /opt/entrypoint.sh
COPY docker/scripts/nginx-healthcheck.sh /opt/healthcheck.sh

EXPOSE 80/tcp
EXPOSE 443/tcp
EXPOSE 8000/tcp
EXPOSE 9200/tcp
HEALTHCHECK --interval=10s --timeout=10s --start-period=5s --retries=3 CMD [ "/opt/healthcheck.sh" ]
ENTRYPOINT [ "/opt/entrypoint.sh" ]