FROM alpine:3.16 as certs
RUN apk --update add ca-certificates

# Use debian 
FROM debian:12-slim

ARG USER_UID=10001
ARG OTEL_BIN=otelcol-dev/otelcol-dev
USER ${USER_UID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 ${OTEL_BIN} /otelcol-contrib
COPY otelcol-contrib.yaml /etc/otelcol-contrib/config.yaml
ENTRYPOINT ["/otelcol-contrib"]
CMD ["--config", "/etc/otelcol-contrib/config.yaml"]
EXPOSE 4317 55678 55679