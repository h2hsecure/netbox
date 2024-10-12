#!/bin/sh
export CONTEXT_PATH="${CONTEXT_PATH:-ddos}"
export CACHE_SIZE="${CACHE_SIZE:-256}"
export PROM_LISTEN="${PROM_LISTEN:-:2112}"
export CLUSTER_ID=1001
export MY_ADDRESS=localhost:12345
export GRPC_SERVER_PORT=45678
export BACKEND_HOST="${BACKEND_HOST:-google.com}"
export BACKEND_PORT="${BACKEND_PORT:-80}"
export DOMAIN="${DOMAIN:-localhost:8080}"
export DOMAIN_PROTO="${DOMAIN_PROTO:-http}"
export DEFAULT_LOCALE="${DEFAULT_LOCALE:-en}"

export CACHE_SOCK=/app/cache.sock
export INTERNAL_SOCK=/app/ddos.sock

chown -R nginx:nginx /app
ls -al /app

umask 0777

# start memcached
/usr/bin/memcached -s $CACHE_SOCK -u nginx -a 0666 -m $CACHE_SIZE -c 1024 -t 4 &

# Start the enforcer
/app/enforce &

# Start the ddos frontend
/app/ddos &

envsubst '$INTERNAL_SOCK $CONTEXT_PATH $BACKEND_HOST $BACKEND_PORT' < /app/nginx.conf > /etc/nginx/nginx.conf

# Start the reverse proxy
/usr/sbin/nginx -g "daemon off;" -c /etc/nginx/nginx.conf &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?