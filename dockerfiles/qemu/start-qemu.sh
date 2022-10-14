#!/bin/bash

set -x 

# TODO fix paths to find the bfp.out
cd /usr/local/bin
./pfloader \
	-cpu host \
	-enable-kvm \
	-display none \
	-serial stdio \
	-nodefaults \
	-m 1024 \
        -device virtio-scsi \
        -object pr-manager-helper,id=helper0,path=/var/run/qemu-pr-helper.sock \
        -blockdev node-name=hd,driver=raw,file.driver=host_device,file.filename=/dev/sdb,file.pr-manager=helper0 \
        -device scsi-block,drive=hd \
	-hda /disk/disk.img
