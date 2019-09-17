FROM ubuntu:18.04 as builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash-completion \
    build-essential \
    ca-certificates \
    curl \
    g++ \
    git \
    jq \
    netcat \
    openssh-client \
    pkg-config \
    python3-pip \
    python3 \
    python \
    python3-dev \
    unzip \
    wget \
    zip \
    zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.12. (121 MiB)
RUN wget -q -Ogo.tgz https://golang.org/dl/go1.12.linux-amd64.tar.gz \
    && echo "750a07fef8579ae4839458701f4df690e0b20b8bcce33b437e4df89c451b6f13 go.tgz" | sha256sum -c - \
    && tar -C /usr/local -xzf go.tgz \
    && rm go.tgz \
    && /usr/local/go/bin/go version

ENV HOME=/ \
    GOPATH=/go \
    PATH=/go/bin:/usr/local/go/bin:/google-cloud-sdk/bin:$PATH \
    GO111MODULE=on \
    GOFLAGS=-mod=vendor

COPY docker/libtensorflow.tar.gz tf.tar.gz
RUN tar -C /usr/local -xzf tf.tar.gz \
    && rm tf.tar.gz \
    && ldconfig

COPY . $GOPATH/src/github.com/unravelin/gotf

WORKDIR $GOPATH/src/github.com/unravelin/gotf

RUN GOOS=linux GOARCH=amd64 go build -o /gotf .

RUN $GOPATH/src/github.com/unravelin/gotf/scripts/copy-libs.sh

FROM scratch

# Copy our static executable.
COPY --from=builder /gotf /gotf
COPY --from=builder /libs /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV LD_LIBRARY_PATH /usr/local/lib:/usr/lib:/lib:/opt/cuda/lib64:/usr/lib64:/lib64

# # No CMD here as it will fail without the models present
COPY model/1/ /model/

EXPOSE 8080

ENTRYPOINT ["./gotf"]
