#!/bin/bash 

set -x

docker rm -f pr-helper
docker rm -f qemu

docker run -ti -d --name pr-helper \
  --pid host \
  --privileged \
  pr-helper

docker run --name qemu --security-opt label=disable \
	--device /dev/sdb:/dev/sdb \
	--device /dev/kvm:/dev/kvm \
	-u root:kvm -td qemu
