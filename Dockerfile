FROM golang:1.14-alpine AS builder

RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache openssl

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

ENV CASSANDRA_URL=cassandra
ENV CASSANDRA_KEYSPACE=gocassandra

WORKDIR /go/src/app/
COPY . /go/src/app/


RUN go get -d -v ./...
RUN go get -d "github.com/google/go-cmp/cmp"
RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app
RUN CGO_ENABLED=0 go test -c -o /bin/handlers_test

EXPOSE 8080

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/app /bin/app
COPY --from=builder /bin/handlers_test /bin/handlers_test
COPY --from=builder /usr/local/bin/dockerize /bin/dockerize