# --------- frontend -----------
FROM node:20.10-alpine as frontendBuilder

WORKDIR /app

ARG proxy

# RUN [ -z "$proxy" ] || sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN npm install -g pnpm@8.14.0
# RUN [ -z "$proxy" ] || pnpm config set registry https://registry.npm.taobao.org

COPY ./web/pnpm-lock.yaml /app/web/pnpm-lock.yaml
COPY ./web/package.json /app/web/package.json

RUN cd /app/web/ && pnpm i

COPY ./web /app/web
RUN cd /app/web/ && pnpm build

COPY ./pal-conf/pnpm-lock.yaml /app/pal-conf/pnpm-lock.yaml
COPY ./pal-conf/package.json /app/pal-conf/package.json

RUN cd /app/pal-conf/ && pnpm i

COPY ./pal-conf /app/pal-conf
RUN cd /app/pal-conf/ && pnpm build

RUN mv /app/pal-conf/dist/assets/* /app/assets
RUN mv /app/pal-conf/dist/index.html /app/pal-conf.html


# --------- sav_cli -----------
FROM python:3.11-alpine as savBuilder

WORKDIR /app

ARG proxy

RUN [ -z "$proxy" ] || sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add build-base

COPY ./module/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt
COPY ./module /app

RUN pyinstaller --onefile sav_cli.py

# --------- backend -----------
FROM golang:1.21-alpine as backendBuilder

ARG proxy
ARG version

WORKDIR /app
ADD . .

COPY --from=frontendBuilder /app/assets /app/assets
COPY --from=frontendBuilder /app/index.html /app/index.html
COPY --from=frontendBuilder /app/pal-conf.html /app/pal-conf.html

RUN if [ ! -z "$proxy" ]; then \
    export GOPROXY=https://goproxy.io,direct && \
    go build -ldflags="-s -w -X 'main.version=${version}'" -o /app/dist/pst main.go; \
    else \
    go build -ldflags="-s -w -X 'main.version=${version}'" -o /app/dist/pst main.go; \
    fi

# --------- runtime -----------
FROM alpine as runtime

WORKDIR /app

ENV SAVE__DECODE_PATH /app/sav_cli

COPY --from=savBuilder /app/dist/sav_cli /app/sav_cli
COPY --from=backendBuilder /app/dist/pst /app/pst

EXPOSE 8080

CMD ["/app/pst"]