FROM registry.access.redhat.com/rhel7-init:7.5

COPY tdservice-*.bin /tmp/tdservice.bin
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN env TDS_NOSETUP=true /tmp/tdservice.bin && rm /tmp/tdservice.bin

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh"]