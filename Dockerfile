FROM golang:latest as builder
WORKDIR /usr/src/techno/forum_app

# Копируем гомод, подтягиваем и кешируем пакеты
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build cmd/application/main.go

FROM ubuntu:18.04 as release

MAINTAINER Danil Rzhevsky

ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER

USER postgres

WORKDIR forum_app
COPY --from=builder /usr/src/techno/forum_app/main .
COPY --from=builder /usr/src/techno/forum_app/db db

RUN service postgresql start &&\
    psql --file=db/db_init.sql &&\
    service postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

RUN cat db/postgresql.conf >> /etc/postgresql/$PGVER/main/postgresql.conf

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

EXPOSE 5000

CMD service postgresql start && ./main



