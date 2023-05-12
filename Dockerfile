FROM --platform=linux/amd64 golang:1.18-alpine

RUN mkdir -p /deploy/app/px
WORKDIR /deploy

COPY app/px/* app/px/
RUN ls -la app/px/
COPY main .

COPY data.json .

EXPOSE 80

CMD [ "/deploy/main" ]