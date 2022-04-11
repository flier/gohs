ARG UBUNTU_VERSION=20.04

FROM ubuntu:${UBUNTU_VERSION}

ARG GO_VERSION=1.17.1
ARG HYPERSCAN_VERSION=5.4.0
ARG PCRE_VERSION=8.45

# Install dependencies

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
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
    zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

# Install golang

RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /go${GO_VERSION}.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf /go${GO_VERSION}.tar.gz && \
    rm /go${GO_VERSION}.tar.gz

ENV PATH=/usr/local/go/bin:${PATH}

# Download Hyperscan

ENV HYPERSCAN_DIR=/hyperscan

RUN wget https://github.com/intel/hyperscan/archive/refs/tags/v${HYPERSCAN_VERSION}.tar.gz -O /hyperscan-${HYPERSCAN_VERSION}.tar.gz && \
    mkdir ${HYPERSCAN_DIR} && tar xf /hyperscan-${HYPERSCAN_VERSION}.tar.gz -C ${HYPERSCAN_DIR} --strip-components=1 && rm /hyperscan-${HYPERSCAN_VERSION}.tar.gz
RUN wget https://sourceforge.net/projects/pcre/files/pcre/${PCRE_VERSION}/pcre-${PCRE_VERSION}.tar.gz/download -O /pcre-${PCRE_VERSION}.tar.gz && \
    mkdir ${HYPERSCAN_DIR}/pcre && tar xf /pcre-${PCRE_VERSION}.tar.gz -C ${HYPERSCAN_DIR}/pcre --strip-components=1 && rm /pcre-${PCRE_VERSION}.tar.gz

# Install Hyperscan

ENV INSTALL_DIR=/usr/local

RUN mkdir ${HYPERSCAN_DIR}/build && cd ${HYPERSCAN_DIR}/build && \
    cmake -G Ninja -DBUILD_STATIC_LIBS=on -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=${INSTALL_DIR} ${HYPERSCAN_DIR} && \
    ninja && ninja install && mv ${HYPERSCAN_DIR}/build/lib/lib*.a ${INSTALL_DIR}/lib/ && cd / && rm -rf ${HYPERSCAN_DIR}

ENV PKG_CONFIG_PATH=${INSTALL_DIR}/lib/pkgconfig

# Add gohs code

ADD . /gohs/

WORKDIR /gohs
ENTRYPOINT ["/usr/local/go/bin/go"]
CMD ["test", "./..."]
