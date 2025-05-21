FROM docker.elastic.co/elasticsearch/elasticsearch:8.15.1

USER root

COPY elasticsearch-analysis-ik-8.15.1.zip /tmp/ik.zip

RUN elasticsearch-plugin install --batch file:///tmp/ik.zip && rm -f /tmp/ik.zip

USER elasticsearch
