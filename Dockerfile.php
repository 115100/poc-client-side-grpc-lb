FROM php:7

RUN apt-get update && \
    apt-get install -y -q git zlib1g-dev && \
    apt-get clean && \
    cd /usr/local/bin && curl -sS https://getcomposer.org/installer | php && \
    cd /usr/local/bin && mv composer.phar composer && \
    pecl install grpc-1.13.0 && \
    docker-php-ext-enable grpc && \
    pecl install protobuf && \
    docker-php-ext-enable protobuf

WORKDIR /opt/greeter-client

COPY composer.json /opt/greeter-client
RUN composer install

COPY . /opt/greeter-client
