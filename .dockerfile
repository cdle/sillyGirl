FROM debian:11

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    sed -i 's/security.debian.org/mirrors.aliyun.com\/debian-security/g' /etc/apt/sources.list

RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    php \
    php-json \
    php-xml \
    php-mbstring \
    php-curl \
    php-gd \
    php-zip \
    php-mysql \
    php-pgsql \
    php-redis \
    php-pear \
    php-dev \
    curl \
    wget \
    git

# 安装gRPC扩展
# RUN pecl install grpc

# 清理缓存和临时文件
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# 下载文件
RUN mkdir -p /usr/local/sillyGirl \
    && ARCH=$(uname -m) \
    && DOWNLOAD_URL="" \
    && if [ "$ARCH" = "x86_64" ]; then \
        DOWNLOAD_URL="/releases/download/main/sillyGirl_linux_amd64"; \
    elif [ "$ARCH" = "aarch64" ]; then \
        DOWNLOAD_URL="/releases/download/main/sillyGirl_linux_arm64"; \
    elif [ "$ARCH" = "armv7l" ]; then \
        DOWNLOAD_URL="/releases/download/main/sillyGirl_linux_armv7"; \
    else \
        echo "Unsupported architecture: $ARCH"; \
        exit 1; \
    fi \
    && curl -L -sSL -f "https://gitee.com/sillybot/sillyGirl$DOWNLOAD_URL" -o /usr/local/sillyGirl/sillyGirl \
    || (echo "Download from original address failed, trying proxy address..." \
    && curl -sSL --connect-timeout 10 -f "https://github.com/cdle/sillyGirl$DOWNLOAD_URL" -o /usr/local/sillyGirl/sillyGirl) \
    && chmod +x /usr/local/sillyGirl/sillyGirl

# 设置工作目录
WORKDIR /usr/local/sillyGirl

ENV PATH="/usr/local/sillyGirl/language/node/yarn/bin:${PATH}"
ENV PATH="/usr/local/sillyGirl/language/node:${PATH}"
ENV SILLYGIRL_DATA_PATH=/usr/local/sillyGirl/

# 指定容器启动时要运行的命令
# CMD ["/usr/local/sillyGirl/sillyGirl", "-t"]

CMD ["/usr/local/sillyGirl/sillyGirl -t"]

# docker build -t sillygirl .
# docker run -d --restart always --name sillygirl sillygirl
