# 使用Alpine镜像作为基础镜像
FROM alpine:latest

# 将APK源更换为国内源（这里使用阿里云的源）
RUN echo "https://mirrors.aliyun.com/alpine/latest-stable/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/latest-stable/community" >> /etc/apk/repositories

# 更新软件包索引并安装必要的工具
RUN apk update && apk add --no-cache \
    python3 \
    py3-pip \
    php \
    php-json \
    php-phar \
    php-openssl \
    php-pdo \
    php-mysqli \
    php-session \
    php-ctype \
    php-tokenizer \
    php-dom \
    php-xml \
    php-xmlwriter \
    php-mbstring \
    php-simplexml \
    php-fileinfo \
    php-opcache \
    php-zlib \
    php-curl \
    php-ftp \
    php-gd \
    php-xmlreader \
    php-pdo_mysql \
    php-pdo_sqlite \
    php-pdo_pgsql \
    php-posix \
    php-sockets \
    php-bcmath \
    php-redis \
    php-pear \
    php-dev \
    php-pear-grpc \
    curl \
    git

# 安装gRPC扩展
RUN pecl install grpc

# 清理缓存和临时文件
RUN rm -rf /tmp/* /var/cache/apk/*

# 下载文件
RUN mkdir -p /usr/local/sillyGirl \
    && ARCH=$(uname -m) \
    && DOWNLOAD_URL="" \
    && if [ "$ARCH" = "x86_64" ]; then \
        DOWNLOAD_URL="https://github.com/cdle/sillyGirl/releases/download/main/sillyGirl_linux_amd64"; \
    elif [ "$ARCH" = "aarch64" ]; then \
        DOWNLOAD_URL="https://github.com/cdle/sillyGirl/releases/download/main/sillyGirl_linux_arm64"; \
    elif [ "$ARCH" = "armv7l" ]; then \
        DOWNLOAD_URL="https://github.com/cdle/sillyGirl/releases/download/main/sillyGirl_linux_armv7"; \
    else \
        echo "Unsupported architecture: $ARCH"; \
        exit 1; \
    fi \
    && curl -sSL --connect-timeout 20 -f "$DOWNLOAD_URL" -o /usr/local/sillyGirl/sillyGirl \
    || (echo "Download from original address failed, trying proxy address..." \
    && curl -sSL --connect-timeout 60 -f https://ghproxy.com/"$DOWNLOAD_URL" -o /usr/local/sillyGirl/sillyGirl) \
    && chmod +x /usr/local/sillyGirl/sillyGirl

# 设置工作目录
WORKDIR /usr/local/sillyGirl





# 指定容器启动时要运行的命令
CMD ["/usr/local/sillyGirl/sillyGirl"]