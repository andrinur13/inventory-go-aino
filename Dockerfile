FROM golang:1.12-alpine

USER root
ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apk add git

RUN mkdir -p /go/src/twc-ota-api

RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/twc-ota-api

COPY ./Gopkg.toml /go/src/twc-ota-api

WORKDIR /go/src/twc-ota-api

RUN chgrp -R 0 ./ && chmod -R g=u ./

RUN dep ensure 

RUN go build

EXPOSE 13003

CMD ["./twc-ota-api"]