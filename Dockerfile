# syntax=docker/dockerfile:1.5-labs

FROM golang:1.19-bullseye as Build


RUN apt install make g++


ADD https://github.com/ggerganov/whisper.cpp.git /whisper
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

FROM debian:bullseye-slim

WORKDIR /subgen

COPY --from=mwader/static-ffmpeg:5.1.2 /ffmpeg /usr/local/bin/
#COPY --from=mwader/static-ffmpeg:5.1.2 /ffprobe /usr/local/bin/
COPY --from=mwader/static-ffmpeg:5.1.2 /versions.json /subgen

COPY --link --from=Build /app/main /subgen
COPY --link --from=Build /whisper/libwhisper.a /subgen


USER 2000:2000

EXPOSE 8080

CMD ["/subgen/main"]