FROM golang:1.12-stretch as build

WORKDIR /app

RUN apt-get update \
 && apt-get install -y ca-certificates  \
 && apt-get install apt-transport-https -y \
 && update-ca-certificates

COPY . .

RUN git config --global http.sslverify false

RUN go build -o bin/go-dmux

FROM golang:1.12-stretch as deploy
WORKDIR /app
COPY --from=build /app/bin/go-dmux .
COPY --from=build /app/docker/config/conf.json ./config/conf.json
COPY --from=build /app/docker/entrypoint.sh .

ENTRYPOINT [ "/app/entrypoint.sh" ]

