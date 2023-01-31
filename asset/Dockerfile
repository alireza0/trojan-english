FROM alpine

LABEL maintainer="Alireza <alireza7@gmail.com>"

ARG TARGETARCH

ARG VERSION_CHECK="https://api.github.com/repos/alireza0/trojan-english/releases/latest"

ARG DOWNLAOD_URL="https://github.com/alireza0/trojan-english/releases/download"

RUN apk add bash openssl unzip curl tzdata && \
    LATEST_VERSION=`curl -H 'Cache-Control: no-cache' -s "$VERSION_CHECK" | grep 'tag_name' | cut -d '"' -f4` && \
    [[ $TARGETARCH =~ "arm64" ]] && ARCH="arm64" || ARCH="amd64" && \
    curl -fL "$DOWNLAOD_URL/$LATEST_VERSION/trojan-linux-$ARCH" -o /usr/local/bin/trojan && \
    chmod +x /usr/local/bin/trojan

EXPOSE 80

CMD ["trojan", "web"]