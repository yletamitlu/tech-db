FROM golang:1.15 AS builder

WORKDIR /build

COPY . .

RUN go build ./cmd/main.go

FROM ubuntu:20.04

EXPOSE 5000
EXPOSE 5432

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install postgresql-12 -y

USER postgres

COPY ./init.sql .

RUN service postgresql start && \
    psql -c "CREATE USER techdbuser WITH superuser login password 'techdb';" && \
    createdb -O techdbuser techdb && \
    psql -d techdb < ./init.sql && \
    service postgresql stop

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /techdb
COPY --from=builder /build/main .

CMD service postgresql start && ./main
