FROM --platform=linux/amd64 golang:1.18-alpine

RUN mkdir -p /deploy/app/px
WORKDIR /deploy

RUN mkdir -p /internal/config

COPY app/px/* app/px/
RUN ls -la app/px/
COPY main .

COPY internal/config/config.yaml internal/config/

COPY data.json .

EXPOSE 80

CMD [ "/deploy/main", "-c", "internal/config/config.yaml"]
