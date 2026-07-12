ARG ALPINE_IMAGE=alpine:3.22@sha256:14358309a308569c32bdc37e2e0e9694be33a9d99e68afb0f5ff33cc1f695dce
ARG NODE_IMAGE=node:20.10-alpine@sha256:9e38d3d4117da74a643f67041c83914480b335c3bd44d37ccf5b5ad86cd715d1
ARG GO_IMAGE=golang:1.21-alpine@sha256:2414035b086e3c42b99654c8b26e6f5b1b1598080d65fd03c7f499552ff4dc94

# --------- web frontend -----------
FROM ${NODE_IMAGE} AS web-builder

WORKDIR /app

RUN npm install -g pnpm@8.14.0

COPY web/pnpm-lock.yaml web/package.json /app/web/
RUN cd /app/web && pnpm install --frozen-lockfile
COPY web /app/web
RUN cd /app/web && pnpm build


# --------- pal-conf frontend -----------
FROM ${NODE_IMAGE} AS pal-conf-builder

WORKDIR /app

RUN npm install -g pnpm@10.12.3

COPY pal-conf/pnpm-lock.yaml pal-conf/package.json /app/pal-conf/
RUN cd /app/pal-conf && pnpm install --frozen-lockfile
COPY pal-conf /app/pal-conf
RUN cd /app/pal-conf && pnpm build

RUN mkdir /app/pal-conf-assets \
    && mv /app/pal-conf/dist/assets/* /app/pal-conf-assets/ \
    && mv /app/pal-conf/dist/index.html /app/pal-conf.html


# --------- sav_cli -----------
FROM ${ALPINE_IMAGE} AS sav-cli-builder

RUN apk add --no-cache python3 py3-pip build-base python3-dev git
RUN python3 -m venv /opt/sav-cli-venv
ENV PATH="/opt/sav-cli-venv/bin:${PATH}"

COPY docker/sav-cli-requirements.lock /tmp/sav-cli-requirements.lock

# PalworldSaveTools commit validated against the real saves under savs/.
ARG PST_TOOLS_REF=8cb429ae3b1460a6a6a0c31c9964ca8cedb65cc5
RUN git clone https://github.com/deafdudecomputers/PalworldSaveTools.git /tmp/pst-tools \
    && git -C /tmp/pst-tools checkout ${PST_TOOLS_REF}

RUN python -m pip install --no-cache-dir -r /tmp/sav-cli-requirements.lock \
    && python -m pip install --no-cache-dir --no-build-isolation --no-deps \
      /tmp/pst-tools/src/palsav/palooz \
    && python -m pip install --no-cache-dir --no-build-isolation --no-deps \
      /tmp/pst-tools/src/palsav \
    && rm -rf /tmp/pst-tools


# --------- map tiles -----------
FROM ${ALPINE_IMAGE} AS map-downloader

WORKDIR /app

RUN apk add --no-cache python3 py3-pip
RUN python3 -m venv /opt/map-venv
COPY docker/map-requirements.lock /tmp/map-requirements.lock
RUN /opt/map-venv/bin/python -m pip install --no-cache-dir -r /tmp/map-requirements.lock
COPY map_down.py /app/map_down.py
RUN /opt/map-venv/bin/python /app/map_down.py \
    && test "$(find /app/map -type f | wc -l | tr -d ' ')" -eq 5461


# --------- backend -----------
FROM ${GO_IMAGE} AS backend-builder

ARG proxy
ARG version

WORKDIR /app
ADD . .

COPY --from=web-builder /app/assets /app/assets
COPY --from=web-builder /app/index.html /app/index.html
COPY --from=pal-conf-builder /app/pal-conf-assets/ /app/assets/
COPY --from=pal-conf-builder /app/pal-conf.html /app/pal-conf.html
COPY --from=map-downloader /app/map /app/map

RUN if [ -n "$proxy" ]; then export GOPROXY=https://goproxy.io,direct; fi \
    && go build -ldflags="-s -w -X 'main.version=${version}'" -o /app/dist/pst main.go


# --------- runtime -----------
FROM ${ALPINE_IMAGE} AS runtime

WORKDIR /app

RUN apk add --no-cache python3 libstdc++
COPY --from=sav-cli-builder /opt/sav-cli-venv /opt/sav-cli-venv
COPY sav_cli/*.py /app/sav_cli_src/
COPY docker/sav-cli-requirements.lock /app/sav_cli_src/requirements.lock
RUN printf '#!/bin/sh\ncd /app/sav_cli_src\nexec /opt/sav-cli-venv/bin/python /app/sav_cli_src/sav_cli.py "$@"\n' \
      > /app/sav_cli \
    && chmod +x /app/sav_cli

ENV SAVE__DECODE_PATH=/app/sav_cli

COPY --from=backend-builder /app/dist/pst /app/pst
COPY LICENSE NOTICE /app/

EXPOSE 8080

CMD ["/app/pst"]
