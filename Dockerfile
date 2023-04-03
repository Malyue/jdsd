FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

FROM ubuntu
RUN apt-get -qq update
RUN apt-get -qq install -y --no-install-recommends ca-certificates curl
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN apt-get install -y tzdata 


COPY --from=builder /build/main /
ENTRYPOINT ["/main"]
