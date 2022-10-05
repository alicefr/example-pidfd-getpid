#!/bin/bash

set -x 

qemu-system-x86_64 \
	-display none \
	-serial stdio \
	-nodefaults \
	-m 1024 \
        -device virtio-scsi \
        -object pr-manager-helper,id=helper0,path=/proxy.sock \
        -blockdev node-name=hd,driver=raw,file.driver=host_device,file.filename=/dev/sdb,file.pr-manager=helper0 \
        -device scsi-block,drive=hd \
	-hda /disk/disk.img
