FROM golang:1.11

WORKDIR /go/src/github.com/DrakeW/redis-cache-proxy

COPY . .

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN dep ensure
RUN go build -o ./server main.go


ENTRYPOINT [ "./server" ]