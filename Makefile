all: images

image-pr-helper: image-pflaume
	docker build -t pr-helper -f dockerfiles/pr-helper/Dockerfile .

image-qemu: image-pflaume image-disk
	docker build -t qemu -f dockerfiles/qemu/Dockerfile .

image-pflaume:
	docker build -t pflaume dockerfiles/pflaume

image-disk:
	docker build --network host -t disk dockerfiles/create-disk

images: image-pr-helper image-qemu
