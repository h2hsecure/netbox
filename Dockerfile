FROM node:18-alpine AS ui

COPY . .
RUN cd ui && npm install && npm run build

FROM golang:1.22-alpine AS builder
RUN apk add --update gcc libc-dev
WORKDIR '/app'
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=ui ui/dist internal/server/ui
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ddos ./cmd/ddos
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o enforce ./cmd/enforce

FROM nginx:1.27-alpine
RUN apk add --update memcached && rm  -rf /tmp/* /var/cache/apk/*

COPY --from=builder /app/ddos /app/ddos
COPY --from=builder /app/enforce /app/enforce

COPY nginx.conf /app/nginx.conf
COPY run.sh /run.sh
RUN chmod +x /run.sh

RUN mkdir -p /logs

EXPOSE 80

HEALTHCHECK --start-period=1m --interval=30s CMD curl --fail http://localhost/ddos/health || exit 1

ENTRYPOINT ["sh", "-c", "/run.sh"]