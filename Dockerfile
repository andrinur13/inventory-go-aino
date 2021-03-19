FROM golang:1.14-alpine

USER root
ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /api-ota
COPY . ./

RUN chgrp -R 0 ./ && chmod -R g=u ./

RUN go build

EXPOSE 13003

CMD ["./api-ota"]