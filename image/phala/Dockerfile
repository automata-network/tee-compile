FROM rust:1.75-bullseye

WORKDIR /workspace
ENV HOME /workspace

RUN rustup toolchain install stable
RUN cargo install cargo-contract

RUN rustup component add rust-src
RUN rustup component add rust-std
RUN rustup component add rustc

#ENV CARGO_TARGET_DIR /workspace/shared/target

RUN cargo contract new test_contract
RUN cd test_contract && cargo contract build && cargo clean

COPY attestable-build-tool /usr/local/sbin

CMD ["bash", "-c", "attestable-build-tool worker -listen vsock://:12345 -dir /workspace"]