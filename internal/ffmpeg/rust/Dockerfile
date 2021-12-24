FROM ubuntu:20.10

RUN apt update

RUN apt-get install -y build-essential curl

RUN apt-get update

RUN apt install -y clang libavcodec-dev libavformat-dev libavutil-dev libavfilter-dev libavdevice-dev pkg-config

RUN curl https://sh.rustup.rs -sSf | bash -s -- -y

RUN echo 'source $HOME/.cargo/env' >> $HOME/.bashrc

RUN mkdir "packets"

WORKDIR "/packets"

COPY Cargo.lock ./
COPY Cargo.toml ./
COPY src ./


