FROM golang:alpine as builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=1

    #设置时区参数
ENV TZ=Asia/Shanghai
RUN sed -i 's!http://dl-cdn.alpinelinux.org/!https://mirrors.ustc.edu.cn/!g' /etc/apk/repositories
RUN apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . ./app
RUN go build -o main .
EXPOSE 3000

CMD ["./main"]