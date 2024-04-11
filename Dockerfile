# builder
FROM golang:alpine as builder

#ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app

RUN go mod download && \
    go build -o go-gin /app/main/

# runner
FROM alpine:latest as runner


ENV TZ=Asia/Shanghai

WORKDIR /

COPY ./entrypoint.sh /

COPY --from=builder /app/go-gin /

RUN apk add --no-cache tzdata ca-certificates jq curl && \
    mkdir -p /var/log/go-gin &&\
    chmod +x /entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
EXPOSE 443