all: local linux

local:
	./build.sh

linux:
	docker run -ti -v $(CURDIR):/go-github/ --workdir /go-github/ qnib/alpn-go-dev ./build.sh

  
