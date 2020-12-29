FROM golang:1.14

ARG app_env
ENV APP_ENV $app_env
ARG action
ENV APP_ACTION $action

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/github.com/zzlpeter/aps-go
COPY . $GOPATH/src/github.com/zzlpeter/aps-go

RUN go build .

EXPOSE 6060

CMD ["./aps-go", "-action=", "${action}"]