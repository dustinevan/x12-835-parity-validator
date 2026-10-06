FROM node:22-bookworm-slim AS ui
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.25-bookworm AS server
# The Amazon RDS certificate authorities, for APP_SERVER_DB_TLS_MODE=verify-full.
# The build fails if the downloaded bundle does not match this checksum; when
# AWS publishes a new bundle, review it and update the checksum.
ARG RDS_CA_BUNDLE_URL=https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem
ARG RDS_CA_BUNDLE_SHA256=e5bb2084ccf45087bda1c9bffdea0eb15ee67f0b91646106e466714f9de3c7e3
RUN mkdir -p /out/data \
 && curl -fsSL "${RDS_CA_BUNDLE_URL}" -o /out/rds-global-bundle.pem \
 && echo "${RDS_CA_BUNDLE_SHA256}  /out/rds-global-bundle.pem" | sha256sum -c -
# Go modules come from vendor/ (make tidy); the build downloads none.
WORKDIR /src
COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY backend/ ./backend/
COPY --from=ui /src/backend/app/site/dist ./backend/app/site/dist
ARG BUILD=container
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.build=${BUILD}" -o /out/validator ./backend/app \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.build=${BUILD}" -o /out/vdb ./backend/app/vdb

# Two programs, one image: /validator (the web service, the default command)
# and /vdb (migrations and the batch pipeline: docker run IMAGE /vdb ...).
# /data is the only directory the programs write.
FROM scratch
COPY --from=server /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=server /out/rds-global-bundle.pem /etc/ssl/rds-global-bundle.pem
COPY --from=server /out/validator /validator
COPY --from=server /out/vdb /vdb
COPY --from=server --chown=65532:65532 /out/data /data
USER 65532:65532
ENV APP_SERVER_WEB_API_HOST=0.0.0.0:8089 \
    APP_SERVER_STORE_ROOT=/data \
    APP_SERVER_DB_TLS_MODE=verify-full \
    APP_SERVER_DB_TLS_ROOT_CERT=/etc/ssl/rds-global-bundle.pem
EXPOSE 8089
STOPSIGNAL SIGTERM
CMD ["/validator"]
