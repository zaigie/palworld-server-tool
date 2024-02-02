# --------- frontend -----------
FROM node:20.10-alpine as frontendBuilder

WORKDIR /app/web

ARG proxy

# RUN [ -z "$proxy" ] || sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN npm install -g pnpm@8.14.0
# RUN [ -z "$proxy" ] || pnpm config set registry https://registry.npm.taobao.org

COPY ./web/pnpm-lock.yaml /app/web/pnpm-lock.yaml
COPY ./web/package.json /app/web/package.json
COPY ./web/* /app/web/

RUN pnpm i && pnpm build

# --------- sav_cli -----------
FROM python:3.11-slim as savBuilder

WORKDIR /app

ARG proxy

RUN [ -z "$proxy" ] || sed -i 's#http://deb.debian.org#https://mirrors.ustc.edu.cn#g' /etc/apt/sources.list.d/debian.sources
RUN apt update && apt install build-essential binutils -y

COPY ./module/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r requirements.txt
COPY ./module /app
RUN pyinstaller --onefile sav_cli.py --name sav_cli

# --------- backend -----------
FROM golang:1.21 as backendBuilder

ARG proxy
ARG version

WORKDIR /app
ADD . .

COPY --from=frontendBuilder /app/dist /app/

RUN if [ ! -z "$proxy" ]; then \
    export GOPROXY=https://goproxy.io,direct && \
    go build -ldflags="-s -w -X 'main.version=${version}'" -o ./dist/pst main.go; \
    else \
    go build -ldflags="-s -w -X 'main.version=${version}'" -o ./dist/pst main.go; \
    fi

# --------- runtime -----------
FROM scratch as runtime

WORKDIR /app

ENV SAVE__DECODE_PATH /app/save_cli

COPY --from=savBuilder /app/dist/sav_cli /app/sav_cli
COPY --from=backendBuilder /app/dist/pst /app/pst

EXPOSE 8080

ENTRYPOINT ["/app/pst"]