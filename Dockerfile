FROM --platform=linux/amd64 golang:1.18-alpine

RUN mkdir -p /deploy/cmd/zk-api-server
COPY build/zk-api-server /deploy/cmd/zk-api-server/
RUN ls -la /deploy/cmd/zk-api-server/

#WORKDIR /deploy


EXPOSE 80

CMD [ "/deploy/cmd/zk-api-server/zk-api-server"]
