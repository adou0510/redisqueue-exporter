FROM alpine

COPY redisqueue-exporter /redisqueue-exporter

CMD ["/redisqueue-exporter"]
