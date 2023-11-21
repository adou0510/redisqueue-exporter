# RedisQueue exporter
[![Docker Pulls](https://img.shields.io/docker/pulls/fjsmlym8177/redisqueue-exporter.svg?maxAge=604800)][hub]

The redisqueue exporter allows probing redis key length,
support type
- list
- zset

## Running this software
### From binaries

### Using the docker image

### Checking the results


## Building the software

### local Build

    GOOS=linux GOARCH=amd64 go build . 

### Building with Docker

After a successful local build:

    docker build -t redisqueue-exporter . 

## Prometheus Configuration

RedisQueue exporter implements the multi-target exporter pattern, so we advice
to read the guide [Understanding and using the multi-target exporter pattern
](https://prometheus.io/docs/guides/multi-target-exporter/) to get the general
idea about the configuration.

The redisqueue exporter needs to be passed the target as a parameter, this can be
done with relabelling.


```yml
scrape_configs:
  - job_name: redisqueue
    metrics_path: /probe	
    static_configs:
    - targets:
      - 192.168.0.104:6379/0/test
      - 192.168.0.104:6379/0/test*
      - pwd@host:6379/0/erp_database_queues*
    relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - target_label: __address__
      replacement: 192.168.0.104:8080
```



[hub]: https://hub.docker.com/repository/docker/fjsmlym8177/redisqueue-exporter/