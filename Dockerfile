FROM golang:1.17.7-alpine3.15 as goBuilder
#RUN apk add --no-cache bash gcc build-base
# 设置阿里云的镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add bash gcc build-base
RUN mkdir -p /server
WORKDIR /server
COPY . .
RUN /bin/bash -c 'if [ ! -e ".env" ]; then  echo "env file not found" ;  exit 1  ; else echo "check env file success" ; exit 0; fi '
RUN go build -o shorturl

FROM alpine:latest
# 设置阿里云的镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone
RUN mkdir -p /app
WORKDIR /app
COPY --from=goBuilder /server/shorturl .
COPY --from=goBuilder /server/.env .
CMD cd /app && ./shorturl serve -d && top