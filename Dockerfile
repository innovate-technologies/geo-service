FROM golang:1.11-alpine as build

RUN apk add --no-cache git gcc musl-dev protobuf

COPY ./ /go/src/github.com/innovate-technologies/geo-service
WORKDIR /go/src/github.com/innovate-technologies/geo-service

RUN go install ./vendor/github.com/golang/protobuf/protoc-gen-go
RUN /usr/bin/protoc -I pb/ pb/geo.proto --go_out=plugins=grpc:pb
RUN cd pb && go get
RUN go build -o geo-server ./cmd/geo-server

FROM alpine

ADD https://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz /var/lib/geodb/GeoLite2-City.mmdb.gz 

RUN apk add --no-cache gzip && \
    cd /var/lib/geodb/ && gzip -d GeoLite2-City.mmdb.gz &&\
    apk del gzip

COPY --from=build /go/src/github.com/innovate-technologies/geo-service/geo-server /usr/local/bin/geo-server

ENV GEOSERVER_GEODB_PATH="/var/lib/geodb/GeoLite2-City.mmdb"

CMD geo-server
