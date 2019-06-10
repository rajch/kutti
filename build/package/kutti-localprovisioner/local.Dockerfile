FROM scratch
WORKDIR /tmp
COPY  out/kutti-localprovisioner .
CMD ["/tmp/kutti-localprovisioner"] 