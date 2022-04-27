# STAGE 1
# Build the executable(s).
FROM golang:1-buster AS stage1

WORKDIR /var/build/go
ADD ./ ./
RUN go build -o /var/build/bin/api ./

#STAGE 2
#Prepare the base image.
FROM debian:latest AS stage2

RUN apt-get update --fix-missing && \
    apt-get install -yqq \
        openssl \
        curl \
        ca-certificates \
        tzdata \
        && \
    apt-get autoclean -yqq && \
    apt-get clean -yqq \

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

FROM stage2 AS stage3

#COPY  templates/ templates/
COPY --from=stage1 /var/build/bin/* /usr/local/bin/
EXPOSE 8080
RUN sed -i "s#http://deb.debian.org#https://deb.debian.org#g" /etc/apt/sources.list
RUN sed -i "s#http://security.debian.org#https://security.debian.org#g" /etc/apt/sources.list

# This will fails with "No system certificates available. Try installing ca-certificates."
RUN apt-get update && apt-get --assume-yes install curl

ENTRYPOINT [ "/usr/local/bin/api" ]