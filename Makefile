all: clean proxy connector

connector: clean-connector
	go build -o connector connector.go

proxy: clean-proxy
	go build -o proxy proxy.go

clean-connector:
	rm -f connector

clean-proxy:
	rm -f proxy

clean: clean-proxy  clean-connector
	rm -rf *.sock

image: proxy connector
	docker build -t getfd .
