# syntax=docker/dockerfile:1.5-labs

FROM golang:1.19-bullseye as Build


RUN apt install make g++


ADD https://github.com/ggerganov/whisper.cpp.git#v1.2.0 /whisper
WORKDIR /whisper

RUN make libwhisper.a

ENV C_INCLUDE_PATH=/whisper
ENV LIBRARY_PATH=/whisper

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o main

FROM ubuntu:22.04

# Install ca-certificates because apparently we might not actually want ssl to work by default?
RUN set -e; \
      export DEBIAN_FRONTEND=noninteractive; \
      apt-get update; \
      apt-get install -y --no-install-recommends ca-certificates && \
      rm -rf /var/cache/apt && \
      apt-get clean && \
      rm -rf /var/lib/apt/lists/*

RUN mkdir "/models" && mkdir /subgen

WORKDIR /subgen

COPY --from=mwader/static-ffmpeg:5.1.2 /ffmpeg /usr/local/bin/
#COPY --from=mwader/static-ffmpeg:5.1.2 /ffprobe /usr/local/bin/
COPY --from=mwader/static-ffmpeg:5.1.2 /versions.json /subgen

COPY --from=Build /app/main /subgen

ENV MODEL_DIR=/models

EXPOSE 8080

CMD ["/subgen/main"]
