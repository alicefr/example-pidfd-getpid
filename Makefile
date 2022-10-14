all: clean proxy connector

connector: clean-connector
	go build -o connector connector.go

proxy: clean-proxy
	go build -o proxy proxy.go

clean-connector:
	rm -f connector

clean: clean-proxy  clean-connector
	rm -rf *.sock

image-pr-helper: connector
	docker build -t pr-helper -f dockerfiles/pr-helper/Dockerfile .

image-qemu: image-pflaume image-disk
	docker build -t qemu -f dockerfiles/qemu/Dockerfile .

image-pflaume:
	docker build -t disk dockerfiles/pflaume

image-disk:
	docker build --network host -t disk dockerfiles/create-disk

images: image-pr-helper image-qemu
