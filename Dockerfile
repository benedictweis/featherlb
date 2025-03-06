FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && \
    apt-get install -y \
    linux-headers-generic \
    clang \
    llvm \
    libbpf-dev \
    ca-certificates \
    curl \
    git \
    bash-completion \
    less \
    vim \
    tree \
    make 

RUN curl -LO https://go.dev/dl/go1.24.1.linux-arm64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.1.linux-arm64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

RUN echo 'PS1="\[\e[32m\]\u@\h:\w\$\[\e[m\] "' >> /root/.bashrc && \
    echo 'alias ls="ls --color=auto"' >> /root/.bashrc && \
    echo 'source /etc/bash_completion' >> /root/.bashrc

WORKDIR /featherlb

COPY . /featherlb

RUN go get ./...

CMD ["bash"]