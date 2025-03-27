ARG GO_VERSION=1.24.1

# ビルド変数で指定したバージョンのgolangを取得
FROM alpine:latest AS latest-go-fetcher
ARG GO_VERSION
ENV GO_VERSION=${GO_VERSION}

RUN apk add --no-cache curl=~8.12.1-r1 tar=~1.35-r2
RUN curl -LO https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
RUN tar -C /tmp -xzf go${GO_VERSION}.linux-amd64.tar.gz

FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm

ARG USERNAME

RUN apt-get update \
    && apt-get install -y \
    vim \
    curl \
    sqlite3 \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts/ \
    && ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

RUN curl -LO https://github.com/BurntSushi/ripgrep/releases/download/14.1.0/ripgrep-14.1.0-x86_64-unknown-linux-musl.tar.gz \
    && tar zxf ripgrep-14.1.0-x86_64-unknown-linux-musl.tar.gz \
    && mv ripgrep-14.1.0-x86_64-unknown-linux-musl/rg /usr/local/bin \
    && rm ripgrep* -rf

RUN if [ -z "${USERNAME}" ] ; then exit 1 ; fi \
    && useradd -u1000 -g1000 -s/bin/bash -b/home/${USERNAME} ${USERNAME} 2> /dev/null || true

RUN install -d -m 0755 -o ${USERNAME} -g ${USERNAME} /workspace

COPY --from=latest-go-fetcher /tmp/go /tmp/go
RUN GOPATH=$(dirname $(dirname $(which go))) \
    && rm ${GOPATH} -rf \
    && mv /tmp/go ${GOPATH}

USER ${USERNAME}

COPY --chown=${USERNAME}:${USERNAME} resources/.vimrc /home/${USERNAME}

RUN go install github.com/spf13/cobra-cli@latest
