build:
	docker build -t register.h2hsecure.com/ddos/worker:latest .

local: build
	docker run --rm -it --name worker -e BACKEND_HOST=pdaccess.com -e BACKEND_PORT=80 -p 8080:80 register.h2hsecure.com/ddos/worker:latest
