FROM ubuntu:22.04

SHELL ["/bin/bash", "-c"]

RUN apt-get update && \
    apt-get install -y gnupg curl && \
    curl -fsSL https://apt.llvm.org/llvm-snapshot.gpg.key | gpg --dearmor -o /usr/share/keyrings/llvm-snapshot.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/llvm-snapshot.gpg] http://apt.llvm.org/jammy/ llvm-toolchain-jammy-20 main" > /etc/apt/sources.list.d/llvm.list

RUN apt-get update && \
    apt-get install -y \
    linux-headers-generic \
    clang \
    llvm-20-dev \
    libclang-20-dev \
    libpolly-20-dev \
    libbpf-dev \
    ca-certificates \
    git \
    bash-completion \
    less \
    vim \
    tree \
    make \
    libssl-dev \
    pkg-config \
    libzstd-dev \
    sudo \
    net-tools \
    iputils-ping

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

RUN source $HOME/.cargo/env && rustup install stable && \
    rustup toolchain install nightly --component rust-src

RUN source $HOME/.cargo/env && cargo install --no-default-features bpf-linker

RUN source $HOME/.cargo/env && cargo install cargo-generate

RUN echo 'PS1="\[\e[32m\]\u@\h:\w\$\[\e[m\] "' >> /root/.bashrc && \
    echo 'alias ls="ls --color=auto"' >> /root/.bashrc && \
    echo 'source /etc/bash_completion' >> /root/.bashrc

WORKDIR /featherlb

COPY . /featherlb

CMD ["bash"]