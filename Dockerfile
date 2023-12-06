FROM --platform=linux/amd64 golang:1.18-alpine

RUN mkdir -p /deploy/cmd/zk-api-server
RUN mkdir -p /deploy/app/px
RUN mkdir -p /internal/config

WORKDIR /deploy

COPY app/px/* app/px/
COPY build/zk-api-server cmd/zk-api-server/
RUN ls -la app/px/

COPY /internal/config/config.yaml internal/config/

EXPOSE 80

CMD [ "cmd/zk-api-server/zk-api-server", "-c", "config/config.yaml"]
