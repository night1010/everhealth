FROM golang:1.18.10-alpine as rest_builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/rest ./cmd/rest/rest.go

FROM postgis/postgis:16-3.4 as postgis_migration

RUN apt update
RUN apt install -y postgis
RUN apt install unzip

WORKDIR /app

COPY ./database /app

RUN #unzip IDN_adm.zip
RUN unzip gadm41_IDN_shp.zip

CMD /app/import.sh

FROM golang:1.18.10-alpine as migration

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD make reset-db

FROM golang:1.18.10-alpine as rest_watch

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD make reload

FROM alpine:3 as rest

WORKDIR /app

COPY --from=rest_builder /app/bin/rest /app/rest
COPY keys.json /app/keys.json
COPY img /app/img

EXPOSE 8080

CMD ./rest
