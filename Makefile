build:
	docker build -t register.h2hsecure.com/ddos/worker:latest .

local: build
	docker run --rm -it --name worker -p 8080:8080 register.h2hsecure.com/ddos/worker:latest
