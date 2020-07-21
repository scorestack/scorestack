FROM docker.elastic.co/elasticsearch/elasticsearch:7.7.1

RUN echo "changeme" | bin/elasticsearch-keystore add -xf bootstrap.password
RUN bin/elasticsearch-users useradd kbn -p changeme -r kibana_system
RUN bin/elasticsearch-users useradd root -p changeme -r superuser