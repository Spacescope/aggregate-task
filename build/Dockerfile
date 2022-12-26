FROM alpine:3.15

RUN mkdir -p /etc/aggregate-task/{bin,conf}
COPY bin/data-api-server /etc/aggregate-task/bin/data-api-server

CMD ["--conf", "/etc/aggregate-task/conf/service.conf"]
ENTRYPOINT ["/etc/aggregate-task/bin/data-api-server"]