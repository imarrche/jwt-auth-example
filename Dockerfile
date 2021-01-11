FROM golang:1.15.5-alpine

RUN apk update && apk add make

WORKDIR ./jwt
COPY . .

RUN make build
