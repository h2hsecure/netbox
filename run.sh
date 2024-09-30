#!/bin/sh


/usr/bin/memcached -p 11211 -u memcached -m 256 -c 1024 -t 4 &

# Start the first process
/app/ddos &

# Start the second process
/usr/sbin/nginx -g "daemon off;" -c /etc/nginx/nginx.conf &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?