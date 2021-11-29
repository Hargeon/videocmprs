# syntax=docker/dockerfile:1

FROM golang:1.16

COPY . /go/src/app

WORKDIR /go/src/app

ENV PORT ":30001"

ENV DB_USER postgres
ENV DB_NAME postgres
ENV DB_SSLMODE disable
ENV DB_HOST localhost
ENV DB_PORT 5433
ENV DB_PASS secret_db
ENV SECRET secret
ENV TOKEN_SECRET secret
ENV AWS_BUCKET_NAME name
ENV AWS_ACCESS_KEY access_key
ENV AWS_SECRET_KEY secret_key
ENV AWS_REGION region
ENV BASE_URL "localhost:3001"
ENV RABBIT_USER guest
ENV RABBIT_PASSWORD guest
ENV RABBIT_HOST localhost
ENV RABBIT_PORT 5672

RUN go build -o videocmprs cmd/videocmprs/main.go

EXPOSE $PORT

CMD ["./videocmprs"]