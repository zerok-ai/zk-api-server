FROM --platform=linux/amd64 golang:1.18-alpine
WORKDIR /zk

COPY build/zk-api-server .
EXPOSE 80

CMD [ "zk-api-server", "-c", "config/config.yaml"]
