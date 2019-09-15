FROM golang:1.13.0-alpine AS BUILDER

WORKDIR /builder

ADD . /builder

RUN go build -o feeds *.go

FROM alpine

WORKDIR /bin

COPY --from=BUILDER  /builder/feeds /bin/feeds

RUN apk add tzdata && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" >  /etc/timezone

ENTRYPOINT [ "/bin/feeds" ,"-p","8888","-d","true"]

EXPOSE 8888:8888
