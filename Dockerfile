FROM golang:1.18.3-bullseye as builder

COPY . /opt
RUN cd /opt && go build -o bin/aggregate-task cmd/aggregate-task/main.go

FROM debian:bullseye
RUN apt update && apt-get install ca-certificates -y
RUN adduser --gecos "Devops Starboard,Github,WorkPhone,HomePhone" --home /app/aggregate-task --disabled-password spacescope
USER spacescope
COPY --from=builder /opt/bin/aggregate-task /app/aggregate-task/aggregate-task

CMD ["--conf", "/app/aggregate-task/service.conf"]
ENTRYPOINT ["/app/aggregate-task/aggregate-task"]
