FROM golang:1.14-alpine3.11 as build
RUN apk --no-cache add tzdata gcc libc-dev git openssh \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && apk del tzdata \
    && rm  -rf /tmp/* /var/cache/apk/*

WORKDIR /var/app/golang
COPY . .
RUN go get -v \
    && go build -o proxy

FROM alpine:3
ENV TZ=Asia/Tokyo
ENV LOG_LEVEL=3
COPY --from=build /var/app/golang/proxy /app/
WORKDIR /app
RUN chmod +x /app/proxy \
    && apk --no-cache add tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && rm  -rf /tmp/* /var/cache/apk/*

EXPOSE 8080

CMD [ "./proxy" ]