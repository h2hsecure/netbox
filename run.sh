#!/bin/sh
export CONTEXT_PATH="${CONTEXT_PATH:-ddos}"
export CACHE_SIZE="${CACHE_SIZE:-256}"
export PROM_LISTEN="${PROM_LISTEN:-:2112}"

#cluster format is host_id:host_name:raft_port:grpc_port
#comma seprated list of hosts
export CLUSTER_STR="${CLUSTER_STR:-1001:ddos1:45678:12345}"
# my format is the same with cluster format
export MY_ADDRESS="${MY_ADDRESS:-1001:ddos1:45678:12345}"

export BACKEND_HOST="${BACKEND_HOST:-google.com}"
export BACKEND_PORT="${BACKEND_PORT:-80}"
export DOMAIN="${DOMAIN:-localhost:8080}"
export DOMAIN_PROTO="${DOMAIN_PROTO:-http}"
export DEFAULT_LOCALE="${DEFAULT_LOCALE:-en}"

export MAX_USER="${MAX_USER:-100}"
export MAX_IP="${MAX_IP:-100}"
export MAX_PATH="${MAX_PATH:-100}"

export ENABLE_SEARCH_ENGINE_BOTS="${ENABLE_SEARCH_ENGINE_BOTS:-true}"

export LOG_DIR=/logs
export CACHE_SOCK=/app/cache.sock
export INTERNAL_SOCK=/app/ddos.sock

chown -R nginx:nginx /app
chown -R nginx:nginx /logs

umask 0777

# start memcached
/usr/bin/memcached -s $CACHE_SOCK -u nginx -a 0666 -m $CACHE_SIZE -c 1024 -t 4 &

# Start the enforcer
/app/enforce &

sleep 3

# Start the ddos frontend
/app/ddos &

envsubst '$INTERNAL_SOCK $CONTEXT_PATH $BACKEND_HOST $BACKEND_PORT' < /app/nginx.conf > /etc/nginx/nginx.conf

# Start the reverse proxy
/usr/sbin/nginx -g "daemon off;" -c /etc/nginx/nginx.conf &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?