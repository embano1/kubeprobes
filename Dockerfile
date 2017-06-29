FROM scratch
MAINTAINER embano1@live.com
COPY probes /probes
ENTRYPOINT ["/probes"]

