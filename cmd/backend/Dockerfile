FROM golang:1.17.7-alpine3.15 AS builder

WORKDIR /go/src/app
COPY . .

RUN go install -v ./...

FROM alpine

WORKDIR /go/src/app

COPY --from=builder /go/bin/backend /go/bin/backend

EXPOSE 8080
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

CMD /wait && /go/bin/backend
