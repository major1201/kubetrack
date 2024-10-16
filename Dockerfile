#FROM golang:1.13.4-alpine as builder
#
#WORKDIR /src
#RUN sed -i s/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g /etc/apk/repositories \
#    && apk --no-cache add ca-certificates git make
#COPY . .
#RUN make linux/amd64

FROM alpine:3.11.3
ENV WEB_LISTEN_ADDR=":80"
CMD [ "kubetrack" ]
EXPOSE 80
RUN sed -i s/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g /etc/apk/repositories \
    && apk --no-cache add ca-certificates
#COPY --from=builder /src/release/kubetrack-* /bin/kubetrack
COPY release/kubetrack-* /bin/kubetrack
