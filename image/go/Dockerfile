FROM golang:1.19.1-bullseye as builder
LABEL maintainer "Automata Team"

WORKDIR /workspace

COPY . /workspace


COPY attestable-build-tool /usr/local/sbin
ENV HOME /workspace

CMD ["bash", "-c", "attestable-build-tool worker -listen vsock://:12345 -dir /workspace"]
