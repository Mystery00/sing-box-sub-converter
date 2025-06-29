FROM golang:1.24-alpine as builder
COPY . /usr/local/go/src/sing-box-sub-converter
WORKDIR /usr/local/go/src/sing-box-sub-converter
RUN GO111MODULE=on go build -o /usr/bin/sing-box-sub-converter sing-box-sub-converter

###
FROM alpine as final
RUN apk add --no-cache ca-certificates && \
        update-ca-certificates
WORKDIR /app
RUN mkdir -p /app/config /app/templates
ENV SUB_CONFIG_HOME /app/config
ENV TEMPLATE_DIR /app/templates
ENTRYPOINT ["/usr/bin/sing-box-sub-converter"]
COPY --from=builder /usr/bin/sing-box-sub-converter /usr/bin/