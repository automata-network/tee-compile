FROM rust:1.75-bullseye

RUN apt update 
#&& apt install -y curl build-essential
#git git-core libssl-dev dpkg-dev autoconf wget ocamlbuild ocaml file pkg-config libtool clang

WORKDIR /workspace
ENV HOME /workspace

#RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
RUN cargo install cargo-contract

RUN rustup component add rust-src --toolchain 1.75.0-x86_64-unknown-linux-gnu
RUN rustup component add cargo --toolchain 1.75.0-x86_64-unknown-linux-gnu
RUN rustup component add clippy --toolchain 1.75.0-x86_64-unknown-linux-gnu
RUN rustup component add rust-std --toolchain 1.75.0-x86_64-unknown-linux-gnu
RUN rustup component add rustc --toolchain 1.75.0-x86_64-unknown-linux-gnu
RUN rustup component add rustfmt --toolchain 1.75.0-x86_64-unknown-linux-gnu

COPY attestable-build-tool /usr/local/sbin

CMD ["bash", "-c", "attestable-build-tool worker -listen vsock://:12345 -dir /workspace"]
