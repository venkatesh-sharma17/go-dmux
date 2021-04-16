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

RUN apt-get update && \
    apt-get install -y procps && \
    apt-get install -y ngrep && \
    apt-get install netcat -y

RUN apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/bin/go-dmux .
COPY --from=build /app/docker/config/conf.json config/conf.json
COPY --from=build /app/docker/entrypoint.sh .

ENV PROJECT "go-dmux"
ENV UID 3119
ENV GID 3000
ENV user "go-dmux"
ENV group "go-dmux"

RUN groupadd -g ${GID} ${group} && \
    useradd -ms /bin/bash --create-home -g ${GID} -u ${UID} ${user}

RUN chown -R ${user}:${group} /var/log && \
    chmod 777 entrypoint.sh && \
    chmod 777 go-dmux

USER ${user}

ENTRYPOINT [ "/app/entrypoint.sh" ]

