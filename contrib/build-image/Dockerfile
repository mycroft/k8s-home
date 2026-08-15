FROM docker.io/golang:1.26.2-trixie

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        ca-certificates=20250419 \
        curl=8.14.1-2+deb13u2 \
        nodejs=20.19.2+dfsg-1+deb13u2 \
        npm=9.2.0~ds1-3 && \
    apt-get clean && \
    rm -fr /var/lib/apt/lists/*

RUN npm install -g cdk8s-cli@2.206.10
