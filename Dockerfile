# syntax=docker/dockerfile:1

ARG UBUNTU_VERSION=22.04

FROM ubuntu:${UBUNTU_VERSION} as build

# Install dependencies

ENV DEBIAN_FRONTEND noninteractive

# hadolint ignore=DL3008
RUN <<EOT bash
    apt-get update
    apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        cmake \
        libboost-dev \
        libbz2-dev \
        libpcap-dev \
        ninja-build \
        pkg-config \
        python2.7 \
        ragel \
        wget \
        zlib1g-dev
    rm -rf /var/lib/apt/lists/*
EOT

# Download Hyperscan

ARG HYPERSCAN_VERSION=5.4.1

ENV HYPERSCAN_DIR=/hyperscan

WORKDIR ${HYPERSCAN_DIR}

ADD https://github.com/intel/hyperscan/archive/refs/tags/v${HYPERSCAN_VERSION}.tar.gz /hyperscan-v${HYPERSCAN_VERSION}.tar.gz
RUN <<EOT bash
    tar xf /hyperscan-v${HYPERSCAN_VERSION}.tar.gz -C ${HYPERSCAN_DIR} --strip-components=1
    rm /hyperscan-v${HYPERSCAN_VERSION}.tar.gz
EOT

ARG PCRE_VERSION=8.45

ADD https://sourceforge.net/projects/pcre/files/pcre/${PCRE_VERSION}/pcre-${PCRE_VERSION}.tar.gz/download /pcre-${PCRE_VERSION}.tar.gz

WORKDIR ${HYPERSCAN_DIR}/pcre

RUN <<EOT bash
    tar xf /pcre-${PCRE_VERSION}.tar.gz -C ${HYPERSCAN_DIR}/pcre --strip-components=1
    rm /pcre-${PCRE_VERSION}.tar.gz
EOT

# Install Hyperscan

ENV INSTALL_DIR=/dist

WORKDIR ${HYPERSCAN_DIR}/build

ARG CMAKE_BUILD_TYPE=RelWithDebInfo

RUN <<EOT bash
    cmake -G Ninja \
        -DBUILD_STATIC_LIBS=ON \
        -DCMAKE_BUILD_TYPE=${CMAKE_BUILD_TYPE} \
        -DCMAKE_INSTALL_PREFIX=${INSTALL_DIR} \
        ..
    ninja
    ninja install
    mv ${HYPERSCAN_DIR}/build/lib/lib*.a ${INSTALL_DIR}/lib/
EOT

FROM ubuntu:${UBUNTU_VERSION}

# Install dependencies

ENV DEBIAN_FRONTEND noninteractive

# hadolint ignore=DL3008
RUN <<EOT bash
    apt-get update
    apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        libpcap-dev \
        pkg-config
    rm -rf /var/lib/apt/lists/*
EOT

# Install golang

ARG GO_VERSION=1.20.3

ADD https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz /

RUN <<EOT bash
    tar -C /usr/local -xzf /go${GO_VERSION}.linux-amd64.tar.gz
    rm /go${GO_VERSION}.linux-amd64.tar.gz
EOT

ENV PATH="/usr/local/go/bin:${PATH}"

ENV INSTALL_DIR=/dist

COPY --from=build ${INSTALL_DIR} ${INSTALL_DIR}

ENV PKG_CONFIG_PATH="${PKG_CONFIG_PATH}:${INSTALL_DIR}/lib/pkgconfig"

# Add gohs code

COPY . /gohs/

WORKDIR /gohs
ENTRYPOINT ["/usr/local/go/bin/go"]
CMD ["test", "-v", "./..."]
