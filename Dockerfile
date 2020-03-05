FROM golang:1.12-alpine3.11

WORKDIR /go/src/app

ENV GO111MODULE=on

RUN apk add --update-cache --virtual .build-deps alpine-sdk git \
    && go get github.com/pilu/fresh \
    && rm -rf /var/cache/apk/*

# set timezone (Alpine)
RUN apk --update-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

EXPOSE 8080

CMD ["fresh"]
