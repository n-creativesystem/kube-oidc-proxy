# FROM golang:1.14-alpine3.11 as build
# RUN apk --no-cache add tzdata gcc libc-dev git openssh \
#     && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
#     && echo "Asia/Tokyo" >  /etc/timezone \
#     && apk del tzdata \
#     && rm  -rf /tmp/* /var/cache/apk/*
FROM golang:1.15-buster as build
ENV TZ=Asia/Tokyo

# RUN apt-get update \
#     && apt-get install -y --no-install-recommends tzdata protobuf-compiler \
#     && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
#     && echo "Asia/Tokyo" >  /etc/timezone \
#     && rm -rf /var/lib/apt/lists/* \
#     && go get google.golang.org/protobuf/cmd/protoc-gen-go \
#         google.golang.org/grpc/cmd/protoc-gen-go-grpc \
#         github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

WORKDIR /var/app/golang
COPY . .
RUN go get -v \
    && go build -o proxy \
    && cd examples/memory \
    && go build -o /var/app/golang/memory

FROM alpine:3
ENV TZ=Asia/Tokyo
ENV LOG_LEVEL=3
WORKDIR /app
COPY --from=build /var/app/golang/proxy /usr/bin/
COPY --from=build /var/app/golang/memory /app/plugin/
COPY docker-entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/docker-entrypoint.sh \
    && chmod +x /usr/bin/proxy \
    && apk --no-cache add tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && rm  -rf /tmp/* /var/cache/apk/*

EXPOSE 8080

CMD [ "./proxy" ]