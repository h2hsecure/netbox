FROM node:18-alpine AS ui

COPY . .
RUN cd ui && npm install && npm run build

FROM golang:1.22-alpine AS builder
RUN apk add --update --no-cache gcc libc-dev
WORKDIR '/app'
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=ui ui/dist internal/server/ui
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ddos ./cmd/ddos
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o enforce ./cmd/enforce

FROM nginx:1.27-alpine
RUN rm /etc/nginx/nginx.conf /etc/nginx/conf.d/default.conf
RUN apk add --no-cache memcached perl

COPY --from=builder /app/ddos /app/ddos
COPY --from=builder /app/enforce /app/enforce
COPY ./tools/mgt /app/mgt

COPY nginx.conf.temp /app/nginx.conf.temp
COPY run.sh /run.sh
RUN chmod +x /run.sh
RUN chmod +x /app/mgt

RUN chown -R nginx:nginx /app
RUN mkdir -p /logs
RUN chown -R nginx:nginx /logs
RUN chown -R nginx:nginx /var/cache/nginx

RUN touch /var/run/nginx.pid && \
        chown -R nginx:nginx /var/run/nginx.pid /run/nginx.pid

RUN touch /app/nginx.conf && \
    chown -R nginx:nginx /app/nginx.conf

EXPOSE 80

HEALTHCHECK --start-period=1m --interval=30s CMD curl --fail http://localhost/ddos/health || exit 1

USER nginx:nginx

ENTRYPOINT ["sh", "-c", "/run.sh"]