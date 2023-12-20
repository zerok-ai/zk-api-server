FROM --platform=linux/amd64 golang:1.18-alpine
RUN mkdir -p /zk
WORKDIR /zk

COPY build/zk-api-server ./
RUN ls -la ./
EXPOSE 80

CMD [ "./zk-api-server", "-c", "./config/config.yaml"]