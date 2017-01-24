FROM alpine:3.4
MAINTAINER Moto Ishizawa "summerwind.jp"

COPY ./h2spec /usr/bin/h2spec

ENTRYPOINT ["h2spec", "--help"]
