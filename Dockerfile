# bulider

FROM golang:alpine AS builder

RUN apk --no-cache --update  add \
    musl-dev \
    util-linux-dev \
    git \
    gcc \ 
    make \
    upx


RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . .

RUN make build

# runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/build/mcrc_tgbot ./

ENTRYPOINT [ "./mcrc_tgbot" ]