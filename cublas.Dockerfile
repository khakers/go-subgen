# syntax=docker/dockerfile:1.5-labs

# using cuda 11.8 due to issues with the whisper.cpp makefile and extra complexities in v1.4.2 with cuda 12
# https://github.com/ggerganov/whisper.cpp/issues/1105

FROM nvidia/cuda:11.8.0-devel-ubuntu22.04 as Build

RUN export DEBIAN_FRONTEND=noninteractive; apt-get update &&  apt-get install make g++ golang -y

# bump: whisper.cpp /whisper.cpp\.git#v([\d.]+)/ https://github.com/ggerganov/whisper.cpp.git|^1
# bump: whisper.cpp link "Release notes" https://github.com/ggerganov/whisper.cpp/releases/tag/v$LATEST
ADD https://github.com/ggerganov/whisper.cpp.git#v1.4.2 /whisper
WORKDIR /whisper

RUN WHISPER_CUBLAS=1 make -j && make WHISPER_CUBLAS=1 libwhisper.a libwhisper.so

ENV C_INCLUDE_PATH="/whisper:/usr/local/cuda/include:/opt/cuda/include:/usr/local/cuda/lib64:/opt/cuda/lib64:/usr/local/cuda-11.8/targets/x86_64-linux/lib/"
ENV LIBRARY_PATH="/whisper"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN export CGO_LDFLAGS="-lcublas -lculibos -lcudart -lcublasLt -lpthread -ldl -lrt -L/usr/local/cuda/lib64 -L/opt/cuda/lib64 -L/usr/local/cuda-11.8/targets/x86_64-linux/lib/"; go build -o main

FROM nvidia/cuda:11.8.0-runtime-ubuntu22.04

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

# bump: static-ffmpeg /static-ffmpeg:([\d.]+)/ docker:mwader/static-ffmpeg|^6
COPY --from=mwader/static-ffmpeg:6.0 /ffmpeg /usr/local/bin/
#COPY --from=mwader/static-ffmpeg:6.0 /ffprobe /usr/local/bin/
COPY --from=mwader/static-ffmpeg:6.0 /versions.json /subgen

COPY --from=build /whisper/libwhisper.so /subgen

ENV LD_LIBRARY_PATH="/subgen"

COPY --from=Build /app/main /subgen

ENV MODEL_DIR=/models

EXPOSE 8095

CMD ["/subgen/main"]