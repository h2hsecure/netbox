FROM node:16-alpine AS ui

COPY . .
RUN cd ui && npm install && npm run build

FROM golang:1.22-alpine AS builder
WORKDIR '/app'
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=ui ui/dist internal/server/ui
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ddos ./cmd/ddos.go 

FROM nginx:1.27-alpine
RUN apk add --update memcached && rm  -rf /tmp/* /var/cache/apk/*

COPY --from=builder /app/ddos /app/ddos

COPY nginx.conf /etc/nginx/nginx.conf
COPY run.sh /run.sh
RUN chmod +x /run.sh

EXPOSE 8080

ENTRYPOINT ["sh", "-c", "/run.sh"]